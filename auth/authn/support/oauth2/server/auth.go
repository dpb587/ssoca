package server

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
)

func (b Service) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	authValue := req.Header.Get("Authorization")
	if authValue == "" {
		return nil, nil
	}

	authValuePieces := strings.SplitN(authValue, " ", 2)
	if len(authValuePieces) != 2 {
		return nil, apierr.NewError(errors.New("invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authValuePieces[0]) != "bearer" {
		return nil, apierr.NewError(errors.New("invalid Authorization method"), http.StatusForbidden, "")
	}

	intTok := selfsignedjwt.NewOriginToken(b.urls.Origin)

	_, err := jwt.ParseWithClaims(
		authValuePieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.JWTSigningMethod {
				return nil, apierr.NewError(errors.New("invalid signing method"), http.StatusForbidden, "")
			}

			return &b.jwtConfig.PrivateKey.PublicKey, nil
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
