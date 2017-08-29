package service

import (
	"net/http"

	"github.com/dpb587/ssoca/service/env/api"
)

//go:generate counterfeiter . Manager
type Manager interface {
	Add(Service)
	Get(string) (Service, error)
	Services() []string
}

//go:generate counterfeiter . Service
type Service interface {
	Type() string
	Version() string
}

type CommandService interface {
	GetCommand() interface{}
	Description() string
}

//go:generate counterfeiter . AuthService
type AuthService interface {
	Service

	AuthRequest(*http.Request) error
	AuthLogin(api.InfoServiceResponse) (interface{}, error)
	AuthLogout() error
}
