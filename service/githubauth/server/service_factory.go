package server

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	oauth2config "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/server/config"
)

type ServiceFactory struct {
	svc.ServiceType

	endpointURL string
	failureURL  string
	successURL  string
}

func NewServiceFactory(endpointURL string, failureURL string, successURL string) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
		failureURL:  failureURL,
		successURL:  successURL,
	}
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWT.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "parsing private key")
	}

	oauthsrv := oauth2server.NewService(
		oauth2config.URLs{
			Origin:      fmt.Sprintf("%s/%s", f.endpointURL, name),
			AuthFailure: f.failureURL,
			AuthSuccess: f.successURL,
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
			PrivateKey:   *privateKey,
			Validity:     *cfg.JWT.Validity,
			ValidityPast: *cfg.JWT.ValidityPast,
		},
	)

	return NewService(name, cfg, oauthsrv), nil
}
