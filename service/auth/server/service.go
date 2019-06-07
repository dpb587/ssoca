package server

import (
	"errors"
	"net/http"

	svc "github.com/dpb587/ssoca/service/auth"
	svcconfig "github.com/dpb587/ssoca/service/auth/server/config"
	svcreq "github.com/dpb587/ssoca/service/auth/server/req"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	svc.Service

	config   svcconfig.Config
	services service.Manager
}

var _ service.Service = Service{}

func NewService(config svcconfig.Config, services service.Manager) Service {
	return Service{
		config:   config,
		services: services,
	}
}

func (s Service) Name() string {
	return "auth"
}

func (s Service) Metadata() interface{} {
	return map[string]interface{}{
		"default_service": s.config.DefaultService,
	}
}

func (s Service) GetRoutes() []req.RouteHandler {
	routes := []req.RouteHandler{
		svcreq.Info{},
	}

	auth, err := s.services.Get(s.config.DefaultService)
	if err != nil {
		// TODO panic?
		return routes
	}

	return append(routes, auth.GetRoutes()...)
}

func (s Service) ParseRequestAuth(r http.Request) (*auth.Token, error) {
	panic(errors.New("TODO"))
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}

func (s Service) getDefaultService() service.AuthService {
	if s.config.DefaultService == "" {
		return nil
	}

	svc, err := s.services.Get(s.config.DefaultService)
	if err != nil {
		panic(err)
	}

	authSvc, ok := svc.(service.AuthService)
	if !ok {
		panic(errors.New("auth service does not implement AuthService"))
	}

	return authSvc
}
