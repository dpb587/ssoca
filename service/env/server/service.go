package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/env"
	svcconfig "github.com/dpb587/ssoca/service/env/config"
	svcreq "github.com/dpb587/ssoca/service/env/server/req"
)

type Service struct {
	svc.Service

	config   svcconfig.Config
	services service.Manager
}

var _ service.Service = Service{}

// @todo config leaking scope
func NewService(config svcconfig.Config, services service.Manager) Service {
	return Service{
		config:   config,
		services: services,
	}
}

func (s Service) Name() string {
	return "env"
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return []req.RouteHandler{
		svcreq.Info{
			Config:   s.config,
			Services: s.services,
		},
	}
}

func (s Service) IsAuthorized(_ http.Request, _ *auth.Token) (bool, error) {
	return true, nil
}
