package service

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service"
)

//go:generate counterfeiter . Factory
type Factory interface {
	Create(service.Type, string, map[string]interface{}) (Service, error)
}

//go:generate counterfeiter . ServiceFactory
type ServiceFactory interface {
	Create(string, map[string]interface{}) (Service, error)
	Type() service.Type
}

//go:generate counterfeiter . Manager
type Manager interface {
	Add(Service)
	Get(string) (Service, error)
	GetAuth() (AuthService, error)
	Services() []string
}

//go:generate counterfeiter . Service
type Service interface {
	Name() string
	Type() service.Type
	Version() string
	Metadata() interface{}
	GetRoutes() []req.RouteHandler
	VerifyAuthorization(http.Request, *auth.Token) error
}

//go:generate counterfeiter . AuthService
type AuthService interface {
	Service

	SupportsRequestAuth(http.Request) (bool, error)
	ParseRequestAuth(http.Request) (*auth.Token, error)
}
