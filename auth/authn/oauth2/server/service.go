package server

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authn/oauth2/server/config"
	svcreq "github.com/dpb587/ssoca/auth/authn/oauth2/server/req"
	"github.com/dpb587/ssoca/auth/authn/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	config            oauth2.Config
	oauthContext      context.Context
	userProfileLoader config.UserProfileLoader

	urls      config.URLs
	jwtConfig config.JWT
}

func NewService(urls config.URLs, config oauth2.Config, oauthContext context.Context, jwtConfig config.JWT, userProfileLoader config.UserProfileLoader) *Service {
	return &Service{
		urls:              urls,
		config:            config,
		oauthContext:      oauthContext,
		jwtConfig:         jwtConfig,
		userProfileLoader: userProfileLoader,
	}
}

func (s Service) VerifyAuthorization(_ http.Request, _ *auth.Token) error {
	return nil
}

func (s Service) SupportsRequestAuth(r http.Request) (bool, error) {
	tokenValue, err := authn.ExtractBearerTokenValue(r)
	if err != nil {
		return false, nil
	} else if tokenValue == "" {
		return false, nil
	}

	claims := selfsignedjwt.NewOriginToken(s.urls.Origin)

	p := &jwt.Parser{}
	_, _, err = p.ParseUnverified(tokenValue, &claims)
	if err != nil {
		return false, nil
	}

	if claims.VerifyOrigin() != nil {
		return false, nil
	}

	return true, nil
}

func (s Service) ParseRequestAuth(r http.Request) (*auth.Token, error) {
	tokenValue, err := authn.ExtractBearerTokenValue(r)
	if err != nil {
		return nil, err
	} else if tokenValue == "" {
		return nil, nil
	}

	intTok := selfsignedjwt.NewOriginToken(s.urls.Origin)

	_, err = jwt.ParseWithClaims(
		tokenValue,
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.JWTSigningMethod {
				return nil, apierr.NewError(errors.New("invalid signing method"), http.StatusForbidden, "")
			}

			return &s.jwtConfig.PrivateKey.PublicKey, nil
		},
	)
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return nil, apierr.NewError(errors.Wrap(err, "parsing claims (ignorable validation error)"), http.StatusUnauthorized, "")
		}

		return nil, apierr.NewError(errors.Wrap(err, "parsing claims"), http.StatusForbidden, "")
	}

	authToken := auth.Token{
		ID:         intTok.ID,
		Groups:     intTok.Groups,
		Attributes: intTok.Attributes,
	}

	return &authToken, nil
}

func (s Service) GetRoutes() []req.RouteHandler {
	return []req.RouteHandler{
		svcreq.Initiate{
			Config: s.config,
		},
		svcreq.Callback{
			URLs:              s.urls,
			UserProfileLoader: s.userProfileLoader,
			Config:            s.config,
			Context:           s.oauthContext,
			JWT:               s.jwtConfig,
		},
	}
}
