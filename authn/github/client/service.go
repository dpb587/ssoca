package client

import (
	"github.com/dpb587/ssoca/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	svc "github.com/dpb587/ssoca/authn/github"
)

type Service struct {
	svc.Service

	runtime   client.Runtime
	cmdRunner boshsys.CmdRunner
}

func NewService(runtime client.Runtime, cmdRunner boshsys.CmdRunner) Service {
	return Service{
		runtime:   runtime,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Description() string {
	return "Authenticate with a GitHub account"
}

func (s Service) GetCommand() interface{} {
	return nil
}
