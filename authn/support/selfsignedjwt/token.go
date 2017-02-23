package selfsignedjwt

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	jwt.StandardClaims

	// @todo extra attributes are signed?
	Username   string            `json:"username,omitempty"`
	Scopes     []string          `json:"scope,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`

	origin string
}

func NewOriginToken(origin string) Token {
	return Token{
		origin: origin,
	}
}

func (t Token) Valid() error {
	// these are not optional for us
	if t.ExpiresAt == 0 {
		return errors.New("Missing exp")
	} else if t.NotBefore == 0 {
		return errors.New("Missing nbf")
	}

	// standard valdiations
	err := t.StandardClaims.Valid()
	if err != nil {
		return err
	}

	// must be self-signed by/for us
	if t.Audience != t.origin {
		return fmt.Errorf("Invalid audience: %s", t.Audience)
	} else if t.Issuer != t.origin {
		return fmt.Errorf("Invalid issuer: %s", t.Issuer)
	}

	return nil
}
