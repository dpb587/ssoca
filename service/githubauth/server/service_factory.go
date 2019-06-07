package server

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	oauth2support "github.com/dpb587/ssoca/auth/authn/support/oauth2"
	oauth2supportconfig "github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/config"
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
	cfg.JWT.Validity = 24 * time.Hour
	cfg.JWT.ValidityPast = 2 * time.Second
	cfg.AuthURL = "https://github.com/login/oauth/authorize"
	cfg.TokenURL = "https://github.com/login/oauth/access_token"

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWT.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "parsing private key")
	}

	backend := oauth2support.NewBackend(
		oauth2supportconfig.URLs{
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
		oauth2supportconfig.JWT{
			PrivateKey:   *privateKey,
			Validity:     cfg.JWT.Validity,
			ValidityPast: cfg.JWT.ValidityPast,
		},
	)

	return NewService(name, cfg, backend), nil
}
