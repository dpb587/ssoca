package oauth2backend

import (
	"errors"
	"net/http"
	"strings"

	"context"

	"golang.org/x/oauth2"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/authn/support/oauth2/config"
	oauth2supportreq "github.com/dpb587/ssoca/authn/support/oauth2/req"
	"github.com/dpb587/ssoca/authn/support/selfsignedjwt"
	"github.com/dpb587/ssoca/server"
	"github.com/dpb587/ssoca/server/service/req"
)

type Backend struct {
	origin       string
	config       oauth2.Config
	oauthContext context.Context

	jwtConfig config.JWT
}

func NewBackend(origin string, config oauth2.Config, oauthContext context.Context, jwtConfig config.JWT) Backend {
	return Backend{
		origin:       origin,
		config:       config,
		oauthContext: oauthContext,

		jwtConfig: jwtConfig,
	}
}

func (b Backend) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	authValue := req.Header.Get("Authorization")
	if authValue == "" {
		authCookie, _ := req.Cookie("Authorization")

		if authCookie == nil {
			return nil, nil
		}

		authValue = authCookie.Value
	}

	authValuePieces := strings.SplitN(authValue, " ", 2)
	if len(authValuePieces) != 2 {
		return nil, server.NewAPIError(errors.New("Invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authValuePieces[0]) != "bearer" {
		return nil, server.NewAPIError(errors.New("Invalid Authorization method"), http.StatusForbidden, "")
	}

	intTok := selfsignedjwt.NewOriginToken(b.origin)

	_, err := jwt.ParseWithClaims(
		authValuePieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.JWTSigningMethod {
				return nil, server.NewAPIError(errors.New("Invalid signing method"), http.StatusForbidden, "")
			}

			return &b.jwtConfig.PrivateKey.PublicKey, nil
		},
	)
	if err != nil {
		return nil, server.NewAPIError(bosherr.WrapError(err, "Parsing claims"), http.StatusForbidden, "")
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
			Origin:            b.origin,
			UserProfileLoader: userProfileLoader,
			Config:            b.config,
			Context:           b.oauthContext,
			JWT:               b.jwtConfig,
		},
	}
}
