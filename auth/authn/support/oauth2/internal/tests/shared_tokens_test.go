package tests_test

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/auth/authn/support/oauth2/internal/tests"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Backend", func() {
	var privateKey *rsa.PrivateKey

	BeforeEach(func() {
		var err error
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(SharedPrivateKey))
		if err != nil {
			panic(err)
		}
	})

	It("one", func() {
		exp, _ := time.Parse(time.RFC3339, "2027-02-14T00:00:00Z")
		iss, _ := time.Parse(time.RFC3339, "2017-02-14T00:00:00Z")
		nbf, _ := time.Parse(time.RFC3339, "2017-02-14T00:00:00Z")
		val := "value1"

		token := jwt.NewWithClaims(config.JWTSigningMethod, selfsignedjwt.Token{
			ID: "fake-user1",
			Groups: []string{
				"scope1",
				"scope2",
			},
			Attributes: map[auth.TokenAttribute]*string{
				auth.TokenAttribute("attr1"): &val,
			},
			StandardClaims: jwt.StandardClaims{
				Audience:  "fake-origin",
				ExpiresAt: exp.Unix(),
				Id:        "fake-uuid",
				IssuedAt:  iss.Unix(),
				Issuer:    "fake-origin",
				NotBefore: nbf.Unix(),
			},
		})

		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			panic(err)
		}

		if tokenString != SharedToken {
			Fail(fmt.Sprintf("Invalid token constant: %s", tokenString))
		}
	})
})
