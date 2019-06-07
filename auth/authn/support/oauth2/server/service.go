package server

import (
	"context"

	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	oauth2supportreq "github.com/dpb587/ssoca/auth/authn/support/oauth2/req"
	"github.com/dpb587/ssoca/server/service/req"
	"golang.org/x/oauth2"
)

type Service struct {
	name string

	config       oauth2.Config
	oauthContext context.Context

	urls              config.URLs
	jwtConfig         config.JWT
	userProfileLoader config.UserProfileLoader
}

func NewService(name string, urls config.URLs, config oauth2.Config, oauthContext context.Context, jwtConfig config.JWT, userProfileLoader config.UserProfileLoader) *Service {
	return &Service{
		name:         name,
		urls:         urls,
		config:       config,
		oauthContext: oauthContext,

		jwtConfig:         jwtConfig,
		userProfileLoader: userProfileLoader,
	}
}

func (s *Service) Name() string {
	return s.name
}

func (s *Service) GetRoutes() []req.RouteHandler {
	return []req.RouteHandler{
		oauth2supportreq.Initiate{
			Config: s.config,
		},
		oauth2supportreq.Callback{
			URLs:              s.urls,
			UserProfileLoader: s.userProfileLoader,
			Config:            s.config,
			Context:           s.oauthContext,
			JWT:               s.jwtConfig,
		},
	}
}
