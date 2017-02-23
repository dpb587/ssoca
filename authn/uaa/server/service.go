package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	svc "github.com/dpb587/ssoca/authn/uaa"
	svc_api "github.com/dpb587/ssoca/authn/uaa/api"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	svc.Service

	name   string
	config Config
}

func NewService(name string, config Config) Service {
	return Service{
		name:   name,
		config: config,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return svc_api.Metadata{
		URL:           s.config.URL,
		CACertificate: s.config.CACertificate,
	}
}

func (s Service) GetRoutes() []req.RouteHandler {
	return nil
}

func (s Service) IsAuthorized(_ http.Request, _ auth.Token) (bool, error) {
	return true, nil
}
