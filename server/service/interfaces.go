package service

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service/req"
)

//go:generate counterfeiter . Factory
type Factory interface {
	Create(string, string, map[string]interface{}) (Service, error)
}

//go:generate counterfeiter . ServiceFactory
type ServiceFactory interface {
	Create(string, map[string]interface{}) (Service, error)
	Type() string
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
	Type() string
	Version() string
	Metadata() interface{}
	GetRoutes() []req.RouteHandler
	IsAuthorized(http.Request, *auth.Token) (bool, error)
}

//go:generate counterfeiter . AuthService
type AuthService interface {
	Service

	ParseRequestAuth(http.Request) (*auth.Token, error)
}
