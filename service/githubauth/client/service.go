package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	oauth "github.com/dpb587/ssoca/auth/authn/support/oauth2/client"
	svc "github.com/dpb587/ssoca/service/githubauth"
)

type Service struct {
	svc.ServiceType
	oauth.Service

	name string
}

var _ service.Service = &Service{}
var _ service.AuthService = &Service{}

func NewService(name string, runtime client.Runtime, cmdRunner boshsys.CmdRunner) Service {
	return Service{
		name:    name,
		Service: oauth.NewService(name, runtime, cmdRunner),
	}
}

func (s Service) Description() string {
	return "Authenticate with a GitHub account"
}

func (s Service) GetCommand() interface{} {
	return nil
}
