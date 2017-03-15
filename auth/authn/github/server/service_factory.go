package server

import (
	"fmt"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	svc "github.com/dpb587/ssoca/auth/authn/github"
	svcconfig "github.com/dpb587/ssoca/auth/authn/github/config"
	oauth2support "github.com/dpb587/ssoca/auth/authn/support/oauth2"
	oauth2supportconfig "github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
)

type ServiceFactory struct {
	endpointURL string
}

func NewServiceFactory(endpointURL string) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
	}
}

func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.JWT.Validity = 24 * time.Hour
	cfg.JWT.ValidityPast = 2 * time.Second
	cfg.AuthURL = "https://github.com/login/oauth/authorize"
	cfg.TokenURL = "https://github.com/login/oauth/access_token"

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWT.PrivateKey))
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing private key")
	}

	backend := oauth2support.NewBackend(
		fmt.Sprintf("%s/%s", f.endpointURL, name),
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
