package server

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	oauth2config "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/googleauth"
	svcconfig "github.com/dpb587/ssoca/service/googleauth/server/config"
	"golang.org/x/oauth2"
)

type Service struct {
	svc.ServiceType

	name   string
	config svcconfig.Config
	oauth  *oauth2server.Service
}

func NewService(name string, rootURL string, config svcconfig.Config) *Service {
	scopes := []string{"https://www.googleapis.com/auth/userinfo.email"}

	if config.Scopes.CloudProject != nil {
		scopes = append(scopes, "https://www.googleapis.com/auth/cloud-platform.read-only")
	}

	oauthsrv := oauth2server.NewService(
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
			Scopes:      scopes,
		},
		oauth2.NoContext,
		oauth2config.JWT{
			PrivateKey:   config.JWT.PrivateKey,
			Validity:     config.JWT.Validity,
			ValidityPast: config.JWT.ValidityPast,
		},
	)

	return &Service{
		name:   name,
		config: config,
		oauth:  oauthsrv,
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
