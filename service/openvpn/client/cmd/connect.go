package cmd

import (
	"fmt"
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Connect struct {
	clientcmd.ServiceCommand

	Exec              string      `long:"exec" description:"Path to the openvpn binary"`
	Reconnect         bool        `long:"reconnect" description:"Reconnect on connection disconnects"`
	StaticCertificate bool        `long:"static-certificate" description:"Write a static certificate in the configuration instead of dynamic renewals"`
	Sudo              bool        `long:"sudo" description:"Execute openvpn with sudo"`
	Args              ConnectArgs `positional-args:"true"`

	FS               boshsys.FileSystem
	CmdRunner        boshsys.CmdRunner
	GetClient        GetClient
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

	client, err := c.GetClient(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	profileManager, err := profile.CreateManagerAndPrivateKey(client, c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting profile manager")
	}

	tmpdir, err := c.FS.TempDir("openvpn")
	if err != nil {
		return bosherr.WrapError(err, "Creating tmpdir")
	}

	defer c.FS.RemoveAll(tmpdir)

	err = c.FS.Chmod(tmpdir, 0700)
	if err != nil {
		return bosherr.WrapError(err, "Chmod'ing tmpdir")
	}

	configPath := fmt.Sprintf("%s/openvpn.ovpn", tmpdir)

	openvpnargs := []string{}

	if c.Sudo {
		openvpnargs = append(openvpnargs, executable)
		executable = "sudo"
	}

	openvpnargs = append(openvpnargs, "--config", configPath)
	openvpnargs = append(openvpnargs, c.Args.Extra...)

	var mgmt management.Client

	if !c.StaticCertificate {
		mgmt = management.NewClient(
			management.NewDefaultHandler(&profileManager),
			"tcp",
			"127.0.0.1:9010",
		)

		mgmt.Start()
		defer mgmt.Stop()
	}

	for {
		profile, err := profileManager.GetProfile()
		if err != nil {
			return bosherr.WrapError(err, "Getting profile")
		}

		if c.StaticCertificate {
			err = c.FS.WriteFileString(configPath, profile.FullConfig())
		} else {
			err = c.FS.WriteFileString(configPath, profile.ManagementConfig(mgmt.ManagementConfigValue()))
		}
		if err != nil {
			return bosherr.WrapError(err, "Writing certificate")
		}

		_, _, _, err = c.CmdRunner.RunComplexCommand(boshsys.Command{
			Name: executable,
			Args: openvpnargs,

			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,

			KeepAttached: true,
		})

		if err != nil {
			return err
		} else if !c.Reconnect {
			return nil
		}
	}
}
