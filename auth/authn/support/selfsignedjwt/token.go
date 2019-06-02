package selfsignedjwt

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"
)

type Token struct {
	jwt.StandardClaims

	ID         string                          `json:"scid,omitempty"`
	Groups     []string                        `json:"scgr,omitempty"`
	Attributes map[auth.TokenAttribute]*string `json:"scat,omitempty"`

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
		return errors.New("missing exp")
	} else if t.NotBefore == 0 {
		return errors.New("missing nbf")
	}

	// standard valdiations
	err := t.StandardClaims.Valid()
	if err != nil {
		return err
	}

	// must be self-signed by/for us
	if t.Audience != t.origin {
		return fmt.Errorf("invalid audience: %s", t.Audience)
	} else if t.Issuer != t.origin {
		return fmt.Errorf("invalid issuer: %s", t.Issuer)
	}

	return nil
}
