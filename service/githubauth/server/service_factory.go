package server

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	oauth2config "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	"github.com/dpb587/ssoca/config"
	serverconfig "github.com/dpb587/ssoca/server/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	endpointURL string
	redirects   serverconfig.ServerRedirectConfig
}

func NewServiceFactory(endpointURL string, redirects serverconfig.ServerRedirectConfig) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
		redirects:   redirects,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	cfg.ApplyRedirectDefaults(f.redirects.AuthSuccess, f.redirects.AuthFailure)

	oauthsrv := oauth2server.NewService(
		oauth2config.URLs{
			Origin:      fmt.Sprintf("%s/%s", f.endpointURL, name),
			AuthFailure: cfg.FailureRedirectURL,
			AuthSuccess: cfg.SuccessRedirectURL,
		},
		oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.AuthURL,
				TokenURL: cfg.TokenURL,
			},
			RedirectURL: fmt.Sprintf("%s/%s/callback", f.endpointURL, name),
			Scopes: []string{
				"read:org",
			},
		},
		oauth2.NoContext,
		oauth2config.JWT{
			PrivateKey:   cfg.JWT.PrivateKey,
			Validity:     *cfg.JWT.Validity,
			ValidityPast: *cfg.JWT.ValidityPast,
		},
	)

	return NewService(name, cfg, oauthsrv), nil
}
