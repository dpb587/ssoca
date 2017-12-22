package cli

import (
	"errors"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/ssh/client"
	"github.com/jessevdk/go-flags"
)

type Agent struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory
	fs             boshsys.FileSystem
	cmdRunner      boshsys.CmdRunner
}

var _ flags.Commander = Agent{}

func (c Agent) Execute(_ []string) error {
	return errors.New("Not yet tested on Windows. Pull requests accepted!")
}
