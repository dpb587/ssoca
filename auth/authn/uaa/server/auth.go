package server

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
	uaainternal "github.com/dpb587/ssoca/auth/authn/uaa/internal"
	apierr "github.com/dpb587/ssoca/server/api/errors"
)

func (s Service) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil
	}

	authHeaderPieces := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderPieces) != 2 {
		return nil, apierr.NewError(errors.New("invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authHeaderPieces[0]) != "bearer" {
		return nil, apierr.NewError(errors.New("invalid Authorization method"), http.StatusForbidden, "")
	}

	intTok := uaainternal.Token{}

	_, err := jwt.ParseWithClaims(
		authHeaderPieces[1],
		&intTok,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method == jwt.SigningMethodNone {
				return nil, apierr.NewError(errors.New("no signing method used"), http.StatusForbidden, "")
			}

			return jwt.ParseRSAPublicKeyFromPEM([]byte(s.config.PublicKey))
		},
	)
	if err != nil {
		return nil, apierr.NewError(errors.Wrap(err, "parsing claims"), http.StatusForbidden, "")
	}

	token := auth.Token{}
	token.ID = intTok.Username
	token.Attributes = map[auth.TokenAttribute]*string{}
	token.Attributes[auth.TokenUsernameAttribute] = &intTok.Username

	for _, scope := range intTok.Scopes {
		token.Groups = append(token.Groups, scope)
	}

	return &token, nil
}
