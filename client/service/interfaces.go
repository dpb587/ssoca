package service

import (
	"net/http"

	"github.com/dpb587/ssoca/service/env/api"
)

//go:generate counterfeiter . Manager
type Manager interface {
	Add(Service)
	AddFactory(ServiceFactory)
	Get(string, string) (Service, error)
}

//go:generate counterfeiter . Service
type Service interface {
	Name() string
	Type() string
	Version() string
}

type ServiceFactory interface {
	New(string) Service
	Type() string
	Version() string
}

//go:generate counterfeiter . AuthService
type AuthService interface {
	Service

	AuthRequest(*http.Request) error
	AuthLogin(api.InfoServiceResponse) (interface{}, error)
	AuthLogout() error
}
