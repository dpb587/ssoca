package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var tunnelblickPreConnectScript = template.Must(template.New("script").Parse(`#!/bin/bash

set -u

REM() { /bin/echo $( date -u +"%Y-%m-%dT%H:%M:%SZ" ) "$@"; }

name="$( basename "$( dirname "$( dirname "$( dirname "$0" )" )" )" )"
shadow="$( dirname "$0" )/config.ovpn"

file="/Users/$USER/Library/Application Support/Tunnelblick/Configurations/$name/Contents/Resources/config.ovpn"

REM "renewing profile"
{{.Exec}} --config "{{.Config}}" --environment "{{.Environment}}" openvpn create-profile --service "{{.Service}}" > "$file.tmp"
exit=$?

if [[ "0" != "$exit" ]]; then
  rm "$file.tmp"

  REM "exiting with failure"
  exit $exit
fi

set -e

REM "patching profile"
echo "remap-usr1 SIGTERM" >> "$file.tmp"

REM "installing profile"
mv -f "$file.tmp" "$file"

REM "installing shadow copy"
cp "$file" "$shadow"

REM done
`))

type CreateTunnelblickProfile struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	SssocaExec string                       `long:"exec-ssoca" description:"Path to the ssoca binary"`
	Name       string                       `long:"name" description:"Specific file name to use for *.tblk"`
	Args       CreateTunnelblickProfileArgs `positional-args:"true"`

	FS        boshsys.FileSystem
	GetClient GetClient
}

var _ flags.Commander = CreateTunnelblickProfile{}

type CreateTunnelblickProfileArgs struct {
	DestinationDir string `positional-arg-name:"DESTINATION-DIR" description:"Directory where the *.tblk profile will be created (default: $PWD)"`
}

func (c CreateTunnelblickProfile) Execute(args []string) error {
	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting config manager")
	}

	exec, err := exec.LookPath(c.SssocaExec)
	if err != nil {
		return bosherr.WrapError(err, "Resolving ssoca executable")
	}

	dir := c.Args.DestinationDir

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return bosherr.WrapError(err, "Getting working directory")
		}
	}

	dir = fmt.Sprintf("%s/%s.tblk", dir, c.Name)

	dirAbs, err := c.FS.ExpandPath(dir)
	if err != nil {
		return bosherr.WrapError(err, "Expanding path")
	}

	err = c.FS.MkdirAll(dirAbs, 0700)
	if err != nil {
		return bosherr.WrapError(err, "Creating target directory")
	}

	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	profile, err := client.BaseProfile()
	if err != nil {
		return bosherr.WrapError(err, "Getting base profile")
	}

	pathConfigOvpn := fmt.Sprintf("%s/config.ovpn", dirAbs)

	err = c.FS.WriteFileString(pathConfigOvpn, profile)
	if err != nil {
		return bosherr.WrapError(err, "Writing config.ovpn")
	}

	err = c.FS.Chmod(pathConfigOvpn, 0400)
	if err != nil {
		return bosherr.WrapError(err, "Chmoding config.ovpn")
	}

	pathPreConnect := fmt.Sprintf("%s/pre-connect.sh", dirAbs)

	var preconnectScriptBuf bytes.Buffer
	err = tunnelblickPreConnectScript.Execute(
		&preconnectScriptBuf,
		struct {
			Exec        string
			Config      string
			Environment string
			Service     string
		}{
			Exec:        exec,
			Config:      configManager.GetSource(),
			Environment: c.Runtime.GetEnvironmentName(),
			Service:     c.ServiceName,
		},
	)
	if err != nil {
		return bosherr.WrapError(err, "Generating Tunnelblick pre-connect.sh")
	}

	preconnectScript := preconnectScriptBuf.String()

	err = c.FS.WriteFileString(pathPreConnect, preconnectScript)
	if err != nil {
		return bosherr.WrapError(err, "Writing pre-connect.sh")
	}

	err = c.FS.Chmod(pathPreConnect, 0500)
	if err != nil {
		return bosherr.WrapError(err, "Chmoding pre-connect.sh")
	}

	return nil
}
