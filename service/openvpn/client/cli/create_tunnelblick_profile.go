package cli

import (
	"fmt"
	"os"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"
)

type CreateTunnelblickProfile struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	SsocaExec string                       `long:"exec-ssoca" description:"Path to the ssoca binary"`
	Name      string                       `long:"name" description:"Specific file name to use for *.tblk"`
	Install   bool                         `long:"install" description:"Install the profile (sudo may prompt for privileges)"`
	Args      createTunnelblickProfileArgs `positional-args:"true"`

	serviceFactory svc.ServiceFactory
	fs             boshsys.FileSystem
	cmdRunner      boshsys.CmdRunner
}

var _ flags.Commander = CreateTunnelblickProfile{}

type createTunnelblickProfileArgs struct {
	DestinationDir string `positional-arg-name:"DESTINATION-DIR" description:"Directory where the *.tblk profile will be created (default: $PWD)"`
}

func (c CreateTunnelblickProfile) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	profile, err := service.CreateTunnelblickProfile(svc.CreateTunnelblickProfileOpts{
		SkipAuthRetry: c.SkipAuthRetry,
		SsocaExec:     c.SsocaExec,
		FileName:      c.Name,
		Directory:     c.Args.DestinationDir,
	})
	if err != nil {
		return err
	}

	if !c.Install {
		return nil
	}

	ui := c.Runtime.GetUI()

	ui.PrintLinef("Attempting to install profile with sudo (your workstation password may be required).")

	defer c.fs.RemoveAll(profile)

	_, _, exit, err := c.cmdRunner.RunComplexCommand(boshsys.Command{
		Name: "sudo",
		Args: []string{
			"--",
			fmt.Sprintf("%s/ssoca-install.sh", profile),
		},

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		KeepAttached: true,
	})
	if exit != 0 && err == nil {
		err = fmt.Errorf("exit code %d", exit)
	}
	if err != nil {
		return errors.Wrap(err, "Installing profile")
	}

	return nil
}
