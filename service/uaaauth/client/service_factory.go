package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/uaaauth"
	"github.com/dpb587/ssoca/service/uaaauth/helper"
)

type ServiceFactory struct {
	svc.ServiceType

	runtime          client.Runtime
	serviceManager   service.Manager
	uaaClientFactory helper.ClientFactory
}

func NewServiceFactory(runtime client.Runtime, serviceManager service.Manager, uaaClientFactory helper.ClientFactory) ServiceFactory {
	return ServiceFactory{
		runtime:          runtime,
		serviceManager:   serviceManager,
		uaaClientFactory: uaaClientFactory,
	}
}

func (sf ServiceFactory) New(name string) service.Service {
	return NewService(name, sf.runtime, sf.serviceManager, sf.uaaClientFactory)
}
