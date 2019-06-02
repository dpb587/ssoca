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
	svc "github.com/dpb587/ssoca/service/googleauth"
	svcconfig "github.com/dpb587/ssoca/service/googleauth/config"
)

type ServiceFactory struct {
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
func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.JWT.Validity = 24 * time.Hour
	cfg.JWT.ValidityPast = 2 * time.Second
	cfg.AuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
	cfg.TokenURL = "https://www.googleapis.com/oauth2/v4/token"

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWT.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "parsing private key")
	}

	scopes := []string{"https://www.googleapis.com/auth/userinfo.email"}

	if cfg.Scopes.CloudProject != nil {
		scopes = append(scopes, "https://www.googleapis.com/auth/cloud-platform.read-only")
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
			Scopes:      scopes,
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
