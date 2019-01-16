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

type CreateLaunchdService struct {
	*clientcmd.ServiceCommand `no-flag:"true"`
	clientcmd.InteractiveAuthCommand

	SsocaExec   string `long:"exec-ssoca" description:"Path to the ssoca binary"`
	Name        string `long:"name" description:"Specific file name to use for *.tblk"`
	OpenvpnExec string `long:"exec-openvpn" description:"Path to the openvpn binary"`
	RunAtLoad   bool   `long:"run-at-load" description:"Run the service at load"`
	LogDir      string `long:"log-dir" description:"Log directory for the service (default: ~/Library/Logs)"`

	Start bool `long:"start" description:"Load and start the service after installation"`

	Args createLaunchdServiceArgs `positional-args:"true"`

	serviceFactory svc.ServiceFactory
	fs             boshsys.FileSystem
	cmdRunner      boshsys.CmdRunner
}

var _ flags.Commander = CreateLaunchdService{}

type createLaunchdServiceArgs struct {
	DestinationDir string `positional-arg-name:"DESTINATION-DIR" description:"Directory where the *.plist service will be created (default: ~/Library/LaunchAgents)"`
}

func (c CreateLaunchdService) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	servicePath, serviceName, err := service.CreateLaunchdService(svc.CreateLaunchdServiceOpts{
		SsocaExec:   c.SsocaExec,
		Name:        c.Name,
		Directory:   c.Args.DestinationDir,
		OpenvpnExec: c.OpenvpnExec,
		RunAtLoad:   c.RunAtLoad,
		LogDir:      c.LogDir,
	})
	if err != nil {
		return err
	}

	if !c.Start {
		return nil
	}

	_, _, exit, err := c.cmdRunner.RunComplexCommand(boshsys.Command{
		Name: "launchctl",
		Args: []string{
			"load",
			servicePath,
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
		return errors.Wrap(err, "loading service")
	}

	_, _, exit, err = c.cmdRunner.RunComplexCommand(boshsys.Command{
		Name: "launchctl",
		Args: []string{
			"start",
			serviceName,
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
		return errors.Wrap(err, "starting service")
	}

	ui := c.Runtime.GetUI()

	ui.PrintLinef("The service '%s' has successfully been started.", serviceName)

	return nil
}
