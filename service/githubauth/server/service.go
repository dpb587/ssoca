package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/server/config"
)

type Service struct {
	svc.ServiceType

	name   string
	config svcconfig.Config
	oauth  *oauth2server.Service
}

func NewService(name string, config svcconfig.Config, oauth *oauth2server.Service) *Service {
	return &Service{
		name:   name,
		config: config,
		oauth:  oauth,
	}
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return s.oauth.GetRoutes(s.OAuthUserProfileLoader)
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
