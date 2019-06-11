package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/uaaauth"
	"github.com/dpb587/ssoca/service/uaaauth/helper"
)

type Service struct {
	svc.ServiceType

	name             string
	runtime          client.Runtime
	uaaClientFactory helper.ClientFactory
}

var _ service.Service = &Service{}
var _ service.AuthService = &Service{}

func NewService(name string, runtime client.Runtime, uaaClientFactory helper.ClientFactory) *Service {
	return &Service{
		name:             name,
		runtime:          runtime,
		uaaClientFactory: uaaClientFactory,
	}
}

func (s Service) Name() string {
	return s.name
}
