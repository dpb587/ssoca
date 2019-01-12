package oauth2backend

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	oauth2supportreq "github.com/dpb587/ssoca/auth/authn/support/oauth2/req"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
)

type Backend struct {
	config       oauth2.Config
	oauthContext context.Context

	urls      config.URLs
	jwtConfig config.JWT
}

func NewBackend(urls config.URLs, config oauth2.Config, oauthContext context.Context, jwtConfig config.JWT) Backend {
	return Backend{
		urls:         urls,
		config:       config,
		oauthContext: oauthContext,

		jwtConfig: jwtConfig,
	}
}

func (b Backend) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	authValue := req.Header.Get("Authorization")
	if authValue == "" {
		return nil, nil
	}

	authValuePieces := strings.SplitN(authValue, " ", 2)
	if len(authValuePieces) != 2 {
		return nil, apierr.NewError(errors.New("Invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authValuePieces[0]) != "bearer" {
		return nil, apierr.NewError(errors.New("Invalid Authorization method"), http.StatusForbidden, "")
	}

	intTok := selfsignedjwt.NewOriginToken(b.urls.Origin)

	_, err := jwt.ParseWithClaims(
		authValuePieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.JWTSigningMethod {
				return nil, apierr.NewError(errors.New("Invalid signing method"), http.StatusForbidden, "")
			}

			return &b.jwtConfig.PrivateKey.PublicKey, nil
		},
	)
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			return nil, apierr.NewError(errors.Wrap(err, "Parsing claims (ignorable validation error)"), http.StatusUnauthorized, "")
		}

		return nil, apierr.NewError(errors.Wrap(err, "Parsing claims"), http.StatusForbidden, "")
	}

	authToken := auth.Token{
		ID:         intTok.ID,
		Groups:     intTok.Groups,
		Attributes: intTok.Attributes,
	}

	return &authToken, nil
}

func (b Backend) GetRoutes(userProfileLoader config.UserProfileLoader) []req.RouteHandler {
	return []req.RouteHandler{
		oauth2supportreq.Initiate{
			Config: b.config,
		},
		oauth2supportreq.Callback{
			URLs:              b.urls,
			UserProfileLoader: userProfileLoader,
			Config:            b.config,
			Context:           b.oauthContext,
			JWT:               b.jwtConfig,
		},
	}
}
