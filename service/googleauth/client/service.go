package client

import (
	"github.com/dpb587/ssoca/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	svc "github.com/dpb587/ssoca/service/googleauth"
)

type Service struct {
	svc.ServiceType

	runtime   client.Runtime
	cmdRunner boshsys.CmdRunner
}

func NewService(runtime client.Runtime, cmdRunner boshsys.CmdRunner) *Service {
	return &Service{
		runtime:   runtime,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Description() string {
	return "Authenticate with a Google account"
}

func (s Service) GetCommand() interface{} {
	return nil
}
