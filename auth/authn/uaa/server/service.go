package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	svc "github.com/dpb587/ssoca/auth/authn/uaa"
	svcapi "github.com/dpb587/ssoca/auth/authn/uaa/api"
	svcconfig "github.com/dpb587/ssoca/auth/authn/uaa/config"
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
	return svcapi.Metadata{
		URL:           s.config.URL,
		CACertificate: s.config.CACertificate,
		ClientID:      s.config.ClientID,
		ClientSecret:  s.config.ClientSecret,
		Prompts:       s.config.Prompts,
	}
}

func (s Service) GetRoutes() []req.RouteHandler {
	return nil
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
