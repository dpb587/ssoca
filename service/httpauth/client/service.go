package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"

	svc "github.com/dpb587/ssoca/service/httpauth"
)

type Service struct {
	svc.ServiceType

	name    string
	runtime client.Runtime
}

var _ service.Service = &Service{}
var _ service.AuthService = &Service{}

func NewService(name string, runtime client.Runtime) *Service {
	return &Service{
		name:    name,
		runtime: runtime,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Description() string {
	return "Authenticate with HTTP username/password"
}

func (s Service) GetCommand() interface{} {
	return nil
}
