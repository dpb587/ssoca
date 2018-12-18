package internal

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	jwt.StandardClaims

	Audience []string `json:"aud"`
	Username string   `json:"user_name"`
	Scopes   []string `json:"scope"`
}
