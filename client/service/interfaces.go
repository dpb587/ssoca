package service

import (
	"net/http"
)

//go:generate counterfeiter . Manager
type Manager interface {
	Add(Service)
	Get(string) (Service, error)
	Services() []string
}

//go:generate counterfeiter . Service
type Service interface {
	// Name() string // TODO
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
	AuthLogin() error
	AuthLogout() error
}
