package server

import (
	"fmt"

	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	oauth2config "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/server/config"
	"golang.org/x/oauth2"
)

type Service struct {
	svc.ServiceType
	*oauth2server.Service

	name   string
	config svcconfig.Config
}

func NewService(name string, rootURL string, config svcconfig.Config) *Service {
	svc := &Service{
		name:   name,
		config: config,
	}

	svc.Service = oauth2server.NewService(
		oauth2config.URLs{
			Origin:      fmt.Sprintf("%s/%s", rootURL, name),
			AuthFailure: config.FailureRedirectURL,
			AuthSuccess: config.SuccessRedirectURL,
		},
		oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthURL,
				TokenURL: config.TokenURL,
			},
			RedirectURL: fmt.Sprintf("%s/%s/callback", rootURL, name),
			Scopes: []string{
				"read:org",
			},
		},
		oauth2.NoContext,
		oauth2config.JWT{
			PrivateKey:   config.JWT.PrivateKey,
			Validity:     config.JWT.Validity,
			ValidityPast: config.JWT.ValidityPast,
		},
		svc.OAuthUserProfileLoader,
	)

	return svc
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}
