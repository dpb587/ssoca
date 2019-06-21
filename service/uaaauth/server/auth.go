package server

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	uaainternal "github.com/dpb587/ssoca/service/uaaauth/internal"
)

func (s Service) SupportsRequestAuth(r http.Request) (bool, error) {
	tokenValue, err := authn.ExtractBearerTokenValue(r)
	if err != nil {
		return false, nil
	} else if tokenValue == "" {
		return false, nil
	}

	claims := uaainternal.Token{}

	p := &jwt.Parser{}
	_, _, err = p.ParseUnverified(tokenValue, &claims)
	if err != nil {
		return false, nil
	}

	if claims.Issuer != s.config.URL {
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

	claims := uaainternal.Token{}

	_, err = jwt.ParseWithClaims(
		tokenValue,
		&claims,
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
	token.ID = claims.Username
	token.Attributes = map[auth.TokenAttribute]*string{}
	token.Attributes[auth.TokenUsernameAttribute] = &claims.Username

	for _, scope := range claims.Scopes {
		token.Groups = append(token.Groups, scope)
	}

	return &token, nil
}
