package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/googleauth"
)

type Service struct {
	svc.ServiceType

	name      string
	runtime   client.Runtime
	cmdRunner boshsys.CmdRunner
}

var _ service.Service = &Service{}
var _ service.AuthService = &Service{}

func NewService(name string, runtime client.Runtime, cmdRunner boshsys.CmdRunner) *Service {
	return &Service{
		name:      name,
		runtime:   runtime,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Name() string {
	return s.name
}
