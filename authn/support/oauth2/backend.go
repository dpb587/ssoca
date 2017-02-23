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

func (b Backend) ParseRequestAuth(req http.Request) (auth.Token, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil
	}

	authHeaderPieces := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderPieces) != 2 {
		return nil, errors.New("Invalid Authorization header format")
	} else if strings.ToLower(authHeaderPieces[0]) != "bearer" {
		return nil, errors.New("Invalid Authorization method")
	}

	intTok := selfsignedjwt.NewOriginToken(b.origin)

	_, err := jwt.ParseWithClaims(
		authHeaderPieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.JWTSigningMethod {
				return nil, errors.New("Invalid signing method")
			}

			return &b.jwtConfig.PrivateKey.PublicKey, nil
		},
	)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing claims")
	}

	attributes := map[string]interface{}{}

	for attrKey, attrValue := range intTok.Attributes {
		attributes[attrKey] = attrValue
	}

	for _, scope := range intTok.Scopes {
		attributes[scope] = true
	}

	return auth.NewSimpleToken(intTok.Username, attributes), nil
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
