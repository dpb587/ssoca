package selfsignedjwt_test

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	. "github.com/dpb587/ssoca/authn/support/selfsignedjwt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	var subject Token

	Describe("Valid", func() {
		BeforeEach(func() {
			now := time.Now()

			subject = NewOriginToken("fake-test")
			subject.StandardClaims = jwt.StandardClaims{
				Audience:  "fake-test",
				ExpiresAt: now.Unix() + 30,
				Issuer:    "fake-test",
				NotBefore: now.Unix() - 30,
			}
		})

		It("is valid", func() {
			err := subject.Valid()

			Expect(err).ToNot(HaveOccurred())
		})

		Context("delegated behavior", func() {
			It("delegates", func() {
				subject.ExpiresAt = subject.StandardClaims.NotBefore

				err := subject.Valid()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("token is expired"))
			})
		})

		Context("errors", func() {
			It("requires exp", func() {
				subject.ExpiresAt = 0

				err := subject.Valid()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Missing exp"))
			})

			It("requires nbf", func() {
				subject.NotBefore = 0

				err := subject.Valid()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Missing nbf"))
			})

			It("requires matching aud", func() {
				subject.Audience = "other"

				err := subject.Valid()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid audience: other"))
			})

			It("requires matching iss", func() {
				subject.Issuer = "other"

				err := subject.Valid()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid issuer: other"))
			})
		})
	})
})
