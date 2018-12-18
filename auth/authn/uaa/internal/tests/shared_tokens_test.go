package tests_test

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uaainternal "github.com/dpb587/ssoca/auth/authn/uaa/internal"
	. "github.com/dpb587/ssoca/auth/authn/uaa/internal/tests"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Tests", func() {
	var privateKey *rsa.PrivateKey

	BeforeEach(func() {
		var err error
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(SharedPrivateKey))
		if err != nil {
			panic(err)
		}
	})

	It("one", func() {
		exp, _ := time.Parse(time.RFC3339, "2027-12-17T00:00:00Z")
		iss, _ := time.Parse(time.RFC3339, "2018-12-17T00:00:00Z")
		nbf, _ := time.Parse(time.RFC3339, "2018-12-17T00:00:00Z")

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, uaainternal.Token{
			Username: "fake-user1",
			Scopes: []string{
				"scope1",
				"scope2",
			},
			Audience: []string{
				"fake-audience1",
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
