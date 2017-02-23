package server

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type internalToken struct {
	jwt.StandardClaims

	Audience []string `json:"aud"`
	Username string   `json:"user_name"`
	Scopes   []string `json:"scope"`
}

func (s Service) ParseRequestAuth(req http.Request) (auth.Token, error) {
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

	intTok := internalToken{}

	_, err := jwt.ParseWithClaims(
		authHeaderPieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method == jwt.SigningMethodNone {
				return nil, errors.New("No signing method used")
			}

			return jwt.ParseRSAPublicKeyFromPEM([]byte(s.config.PublicKey))
		},
	)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing claims")
	}

	attributes := map[string]interface{}{}

	for _, scope := range intTok.Scopes {
		attributes[scope] = true
	}

	return auth.NewSimpleToken(intTok.Username, attributes), nil
}
