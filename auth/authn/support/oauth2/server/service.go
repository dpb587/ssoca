package server

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	svcreq "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/req"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	config       oauth2.Config
	oauthContext context.Context

	urls      config.URLs
	jwtConfig config.JWT
}

func NewService(urls config.URLs, config oauth2.Config, oauthContext context.Context, jwtConfig config.JWT) *Service {
	return &Service{
		urls:         urls,
		config:       config,
		oauthContext: oauthContext,

		jwtConfig: jwtConfig,
	}
}

func (s Service) ParseRequestAuth(r http.Request) (*auth.Token, error) {
	authValue := r.Header.Get("Authorization")
	if authValue == "" {
		return nil, nil
	}

	authValuePieces := strings.SplitN(authValue, " ", 2)
	if len(authValuePieces) != 2 {
		return nil, apierr.NewError(errors.New("invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authValuePieces[0]) != "bearer" {
		return nil, apierr.NewError(errors.New("invalid Authorization method"), http.StatusForbidden, "")
	}

	intTok := selfsignedjwt.NewOriginToken(s.urls.Origin)

	_, err := jwt.ParseWithClaims(
		authValuePieces[1],
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

func (s Service) GetRoutes(userProfileLoader config.UserProfileLoader) []req.RouteHandler {
	return []req.RouteHandler{
		svcreq.Initiate{
			Config: s.config,
		},
		svcreq.Callback{
			URLs:              s.urls,
			UserProfileLoader: userProfileLoader,
			Config:            s.config,
			Context:           s.oauthContext,
			JWT:               s.jwtConfig,
		},
	}
}
