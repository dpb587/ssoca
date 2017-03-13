package cmd

import (
	"fmt"
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Connect struct {
	clientcmd.ServiceCommand

	Exec string `long:"exec" description:"Path to the openvpn binary"`
	// Reconnect bool `long:"reconnect" description:"Reconnect on connection disconnects"`
	Sudo bool        `long:"sudo" description:"Execute openvpn with sudo"`
	Args ConnectArgs `positional-args:"true"`

	FS               boshsys.FileSystem
	CmdRunner        boshsys.CmdRunner
	CreateProfile    CreateUserProfile
	ExecutableFinder client.ExecutableFinder
}

var _ flags.Commander = Connect{}

type ConnectArgs struct {
	Extra []string `positional-arg-name:"EXTRA" description:"Additional arguments to pass to openvpn"`
}

func (c Connect) Execute(args []string) error {
	var executable string

	if c.Exec != "" {
		executable = c.Exec
	} else {
		var err error

		executable, err = c.ExecutableFinder.Find()
		if err != nil {
			return bosherr.WrapError(err, "Finding executable")
		}
	}

	profile, err := c.CreateProfile(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Creating profile")
	}

	tmpdir, err := c.FS.TempDir("openvpn")
	if err != nil {
		return bosherr.WrapError(err, "Creating openvpn tmpdir")
	}

	defer c.FS.RemoveAll(tmpdir)

	configPath := fmt.Sprintf("%s/openvpn.ovpn", tmpdir)

	err = c.FS.WriteFileString(configPath, profile)
	if err != nil {
		return bosherr.WrapError(err, "Writing certificate")
	}

	openvpnargs := []string{}

	if c.Sudo {
		openvpnargs = append(openvpnargs, executable)
		executable = "sudo"
	}

	openvpnargs = append(openvpnargs, "--config", configPath)
	openvpnargs = append(openvpnargs, c.Args.Extra...)

	_, _, _, err = c.CmdRunner.RunComplexCommand(boshsys.Command{
		Name: executable,
		Args: openvpnargs,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		KeepAttached: true,
	})

	return nil
}
