package server_test

import (
	"fmt"
	"net/http"

	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth"
	internaltests "github.com/dpb587/ssoca/auth/authn/oauth2/internal/tests"
	. "github.com/dpb587/ssoca/auth/authn/oauth2/server"
	"github.com/dpb587/ssoca/auth/authn/oauth2/server/config"
	apierr "github.com/dpb587/ssoca/server/api/errors"

	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var subject *Service
	var privateKey *rsa.PrivateKey

	BeforeEach(func() {
		var err error
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(internaltests.SharedPrivateKey))
		if err != nil {
			panic(err)
		}
	})

	Describe("ParseRequestAuth + SupportsRequestAuth", func() {
		BeforeEach(func() {
			subject = NewService(
				config.URLs{Origin: "fake-origin"},
				oauth2.Config{},
				oauth2.NoContext,
				config.JWT{
					PrivateKey: &config.PrivateKey{
						PrivateKey: privateKey,
					},
				},
				nil,
			)
		})

		Context("when valid token is used", func() {
			var tokenString = internaltests.SharedToken
			var req http.Request

			BeforeEach(func() {
				req = http.Request{
					Header: http.Header{
						"Authorization": []string{
							fmt.Sprintf("bearer %s", tokenString),
						},
					},
				}
			})

			It("works", func() {
				supports, err := subject.SupportsRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(supports).To(BeTrue())

				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token.ID).To(Equal("fake-user1"))
			})

			It("sets attributes", func() {
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(*token.Attributes[auth.TokenAttribute("attr1")]).To(Equal("value1"))
			})

			It("sets groups", func() {
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token.Groups).To(ContainElement("scope1"))
				Expect(token.Groups).To(ContainElement("scope2"))
			})
		})

		Context("invalid token data", func() {
			Context("with invalid algorithm", func() {
				var tokenString = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJmYWtlLW9yaWdpbiIsImV4cCI6MTgwMjY0MTA2NCwianRpIjoic29tZS11dWlkIiwiaWF0IjoxNDg3MjgxMDY0LCJpc3MiOiJmYWtlLW9yaWdpbiIsIm5iZiI6MTE3MTkyMTA2NCwic3ViIjoidXNlZnVsIHN1YmplY3QiLCJ1c2VybmFtZSI6InRlc3QtdXNlciIsInNjb3BlIjpbInNjb3BlMSIsInNjb3BlMiJdLCJhdHRyaWJ1dGVzIjp7ImF0dHIxIjoidmFsdWUxIiwiYXR0cjIiOiIifX0.K3kRT2mbAtR_I7A_7zpRHnW_KggBHbCEqjYjhAt6mzqWcYQuEhVfbQerErqD3coyxoC4zIKsdnubtgl6ijWTsw"

				It("errors", func() {
					req := http.Request{
						Header: http.Header{
							"Authorization": []string{
								fmt.Sprintf("bearer %s", tokenString),
							},
						},
					}

					supports, err := subject.SupportsRequestAuth(req)

					Expect(err).ToNot(HaveOccurred())
					Expect(supports).To(BeTrue())

					token, err := subject.ParseRequestAuth(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("invalid signing method"))

					errapi, errok := err.(apierr.Error)
					Expect(errok).To(BeTrue())
					Expect(errapi.Status).To(Equal(http.StatusUnauthorized))

					Expect(token).To(BeNil())
				})
			})

			Context("with invalid algorithm", func() {
				var tokenString = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJmYWtlLW9yaWdpbiIsImV4cCI6MTgwMjY0MTA2NCwianRpIjoic29tZS11dWlkIiwiaWF0IjoxNDg3MjgxMDY0LCJpc3MiOiJmYWtlLW9yaWdpbiIsIm5iZiI6MTE3MTkyMTA2NCwic3ViIjoidXNlZnVsIHN1YmplY3QiLCJ1c2VybmFtZSI6InRlc3QtdXNlciIsInNjb3BlIjpbInNjb3BlMSIsInNjb3BlMiJdLCJhdHRyaWJ1dGVzIjp7ImF0dHIxIjoidmFsdWUxIiwiYXR0cjIiOiIifX0.K3kRT2mbAtR_I7A_7zpRHnW_KggBHbCEqjYjhAt6mzqWcYQuEhVfbQerErqD3coyxoC4zIKsdnubtgl6ijWT"

				It("errors", func() {
					req := http.Request{
						Header: http.Header{
							"Authorization": []string{
								fmt.Sprintf("bearer %s", tokenString),
							},
						},
					}

					supports, err := subject.SupportsRequestAuth(req)

					Expect(err).ToNot(HaveOccurred())
					Expect(supports).To(BeTrue())

					token, err := subject.ParseRequestAuth(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("parsing claims (ignorable validation error)"))

					errapi, errok := err.(apierr.Error)
					Expect(errok).To(BeTrue())
					Expect(errapi.Status).To(Equal(http.StatusUnauthorized))

					Expect(token).To(BeNil())
				})
			})

			Context("with mismatching issuer", func() {
				BeforeEach(func() {
					subject = NewService(
						config.URLs{Origin: "fake-fake-origin"},
						oauth2.Config{},
						oauth2.NoContext,
						config.JWT{
							PrivateKey: &config.PrivateKey{
								PrivateKey: privateKey,
							},
						},
						nil,
					)
				})

				It("errors", func() {
					req := http.Request{
						Header: http.Header{
							"Authorization": []string{
								fmt.Sprintf("bearer %s", internaltests.SharedToken),
							},
						},
					}

					supports, err := subject.SupportsRequestAuth(req)

					Expect(err).ToNot(HaveOccurred())
					Expect(supports).To(BeFalse())

					_, err = subject.ParseRequestAuth(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("parsing claims (ignorable validation error)"))
				})
			})
		})

		Context("missing Authorization header", func() {
			It("bypasses", func() {
				req := http.Request{}

				supports, err := subject.SupportsRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(supports).To(BeFalse())

				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(BeNil())
			})
		})

		Context("invalid Authorization header", func() {
			It("requires two segments", func() {
				req := http.Request{
					Header: http.Header{
						"Authorization": []string{
							"bearer",
						},
					},
				}

				supports, err := subject.SupportsRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(supports).To(BeFalse())

				token, err := subject.ParseRequestAuth(req)

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})

			It("requires bearer type", func() {
				req := http.Request{
					Header: http.Header{
						"Authorization": []string{
							"basic something",
						},
					},
				}

				supports, err := subject.SupportsRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(supports).To(BeFalse())

				token, err := subject.ParseRequestAuth(req)

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})
		})
	})
})
