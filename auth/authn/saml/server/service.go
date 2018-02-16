package server

import (
	"net/http"

	"github.com/crewjam/saml"
	"github.com/dpb587/ssoca/auth"
	svc "github.com/dpb587/ssoca/auth/authn/saml"
	svcconfig "github.com/dpb587/ssoca/auth/authn/saml/config"
	svcreq "github.com/dpb587/ssoca/auth/authn/saml/server/req"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	svc.Service

	name   string
	config svcconfig.Config
	idp    *saml.ServiceProvider
}

func NewService(name string, config svcconfig.Config, idp *saml.ServiceProvider) Service {
	return Service{
		name:   name,
		config: config,
		idp:    idp,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return []req.RouteHandler{
		svcreq.Initiate{
			Config: s.config,
			IDP:    s.idp,
		},
		svcreq.Callback{
			Config: s.config,
			IDP:    s.idp,
		},
	}
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
