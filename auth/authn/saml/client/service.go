package client

import (
	"github.com/dpb587/ssoca/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	svc "github.com/dpb587/ssoca/auth/authn/saml"
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
	return "Authenticate with a SAML Identity Provider"
}

func (s Service) GetCommand() interface{} {
	return nil
}
