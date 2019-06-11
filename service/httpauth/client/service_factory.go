package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/httpauth"
)

type ServiceFactory struct {
	svc.ServiceType

	runtime client.Runtime
}

func NewServiceFactory(runtime client.Runtime) ServiceFactory {
	return ServiceFactory{
		runtime: runtime,
	}
}

func (sf ServiceFactory) New(name string) service.Service {
	return NewService(name, sf.runtime)
}
