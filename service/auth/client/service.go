package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"

	svc "github.com/dpb587/ssoca/service/auth"
	svchttpclient "github.com/dpb587/ssoca/service/auth/httpclient"
)

type Service struct {
	svc.ServiceType

	runtime        client.Runtime
	serviceManager service.Manager
}

var _ service.Service = Service{}

func NewService(runtime client.Runtime, serviceManager service.Manager) *Service {
	return &Service{
		runtime:        runtime,
		serviceManager: serviceManager,
	}
}

func (s Service) Name() string {
	return "auth"
}

func (s Service) Description() string {
	return "Manage authentication"
}

func (s Service) GetServiceManager() service.Manager {
	return s.serviceManager
}

func (s Service) GetClient() (svchttpclient.Client, error) {
	client, err := s.runtime.GetClient()
	if err != nil {
		return nil, err
	}

	return svchttpclient.New(client)
}
