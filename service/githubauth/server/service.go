package server

import (
	"net/http"
	"path"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"
	oauthconfig "github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	oauth "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	svc "github.com/dpb587/ssoca/service/githubauth"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Service struct {
	svc.ServiceType
	*oauth.Service

	name   string
	config svcconfig.Config
}

func NewService(name string, endpointURL string, config svcconfig.Config) *Service {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.JWT.PrivateKey))
	if err != nil {
		panic(errors.Wrap(err, "parsing private key")) // TODO native YAML unmarshal type to capture error earlier
	}

	service := &Service{
		name:   name,
		config: config,
	}

	service.Service = oauth.NewService(
		name,
		oauthconfig.URLs{
			Origin:      path.Join(endpointURL, name),
			AuthFailure: config.RedirectFailure,
			AuthSuccess: config.RedirectSuccess,
		},
		oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthURL,
				TokenURL: config.TokenURL,
			},
			RedirectURL: path.Join(endpointURL, name, "callback"),
			Scopes: []string{
				"read:org",
			},
		},
		oauth2.NoContext,
		oauthconfig.JWT{
			PrivateKey:   *privateKey,
			Validity:     *config.JWT.Validity,
			ValidityPast: *config.JWT.ValidityPast,
		},
		service.BuildAuthToken,
	)

	return service
}

func (s Service) Name() string {
	return s.name
}

func (s Service) Metadata() interface{} {
	return nil
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}
