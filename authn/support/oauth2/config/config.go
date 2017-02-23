package config

import (
	"crypto/rsa"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWTSigningMethod = jwt.SigningMethodRS256

const CookieStateName = "ssoca_oauth_state"
const CookieClientPortName = "ssoca_oauth_clientport"

type UserProfileLoader func(*http.Client) (UserProfile, error)

type UserProfile struct {
	Username   string
	Scopes     []string
	Attributes map[string]string
}

type JWT struct {
	PrivateKey   rsa.PrivateKey
	Validity     time.Duration
	ValidityPast time.Duration
}
