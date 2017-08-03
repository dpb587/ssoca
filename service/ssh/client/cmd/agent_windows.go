package cmd

import (
	"errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/jessevdk/go-flags"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Agent struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	Foreground bool   `long:"foreground" description:"Stay in foreground"`
	Socket     string `long:"socket" description:"Socket path (ensure the directory has restricted permissions)"`

	GetClient GetClient
	CmdRunner boshsys.CmdRunner
	FS        boshsys.FileSystem
}

var _ flags.Commander = Agent{}

func (c Agent) Execute(args []string) error {
	return errors.New("Not yet tested on Windows. Pull requests accepted!")
}
