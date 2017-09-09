package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	svc "github.com/dpb587/ssoca/auth/authn/http"
	svcconfig "github.com/dpb587/ssoca/auth/authn/http/config"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	svc.Service

	name   string
	config svcconfig.Config
}

func NewService(name string, config svcconfig.Config) Service {
	return Service{
		name:   name,
		config: config,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return nil
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
