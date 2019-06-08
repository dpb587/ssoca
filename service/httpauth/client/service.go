package client

import (
	"github.com/dpb587/ssoca/client"

	svc "github.com/dpb587/ssoca/service/httpauth"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.ServiceType

	runtime client.Runtime
	fs      boshsys.FileSystem
}

func NewService(runtime client.Runtime) *Service {
	return &Service{
		runtime: runtime,
	}
}

func (s Service) Description() string {
	return "Authenticate with HTTP username/password"
}

func (s Service) GetCommand() interface{} {
	return nil
}
