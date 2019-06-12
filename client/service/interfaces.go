package service

import (
	"net/http"

	"github.com/dpb587/ssoca/service"
	"github.com/dpb587/ssoca/service/env/api"
)

//go:generate counterfeiter . Manager
type Manager interface {
	Add(Service)
	AddFactory(ServiceFactory)
	Get(service.Type, string) (Service, error)
}

//go:generate counterfeiter . Service
type Service interface {
	Name() string
	Type() service.Type
	Version() string
}

type ServiceFactory interface {
	New(string) Service
	Type() service.Type
	Version() string
}

//go:generate counterfeiter . AuthService
type AuthService interface {
	Service

	AuthRequest(*http.Request) error
	AuthLogin(api.InfoServiceResponse) (interface{}, error)
	AuthLogout() error
}
