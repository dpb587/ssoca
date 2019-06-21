package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	oauthsvc "github.com/dpb587/ssoca/auth/authn/oauth2/client"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/googleauth"
)

type Service struct {
	svc.ServiceType
	*oauthsvc.AuthService

	name      string
	runtime   client.Runtime
	cmdRunner boshsys.CmdRunner
}

var _ service.Service = &Service{}
var _ service.AuthService = &Service{}

func NewService(name string, runtime client.Runtime, cmdRunner boshsys.CmdRunner) *Service {
	srv := &Service{
		name:      name,
		runtime:   runtime,
		cmdRunner: cmdRunner,
	}

	srv.AuthService = oauthsvc.NewAuthService(name, srv.Type(), runtime, cmdRunner)

	return srv
}

func (s Service) Name() string {
	return s.name
}
