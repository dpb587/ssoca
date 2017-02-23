package oauth2backend_test

import (
	"fmt"
	"net/http"

	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/dpb587/ssoca/authn/support/oauth2"
	oauth2supportconfig "github.com/dpb587/ssoca/authn/support/oauth2/config"

	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend", func() {
	var subject Backend
	var privateKey *rsa.PrivateKey

	BeforeEach(func() {
		var err error
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALx3tjO/le0iinRqpdMa3/+4lbYuPQOsJROsDo9D4gN3rvnqOiuI
MIydz3VZ+ErHf5+LCQa4elzMK6mGpIVrewcCAwEAAQJBAJfHZPHZ6fkWpyBIPxF7
BEhiNBeKt1J80UM9fmA8UAlvd3m003ANZAIe6/GEYCTVp49jJ0e66Fwhd+LOsFN2
cPkCIQDKSnqE8QKg45S/QosRuCpM2EKwscp0855wdRAyQFYeWwIhAO6BrRRbAriQ
8tafNqIeKBaKBLZlen0ataAPiKAEqo3FAiAR1bMrmVwT9zycCC/ephAEqmRm06X3
3aqwW4HMDGQLVQIgUZL6rp6eJKA23l8gIXys+2CDUhsNNOLAwhjuAsT1zH0CIQC/
Y5S1jt/4jLs8Xfh6zPs15EuvQKpUGYKAIHofIp0T3Q==
-----END RSA PRIVATE KEY-----`))
		if err != nil {
			panic(err)
		}
	})

	Describe("ParseRequestAuth", func() {
		BeforeEach(func() {
			subject = NewBackend(
				"fake-origin",
				oauth2.Config{},
				oauth2.NoContext,
				oauth2supportconfig.JWT{
					PrivateKey: *privateKey,
				},
			)
		})

		Context("when valid token is used", func() {
			var tokenString = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJmYWtlLW9yaWdpbiIsImV4cCI6MTgwMjY0MTA2NCwianRpIjoic29tZS11dWlkIiwiaWF0IjoxNDg3MjgxMDY0LCJpc3MiOiJmYWtlLW9yaWdpbiIsIm5iZiI6MTE3MTkyMTA2NCwic3ViIjoidXNlZnVsIHN1YmplY3QiLCJ1c2VybmFtZSI6InRlc3QtdXNlciIsInNjb3BlIjpbInNjb3BlMSIsInNjb3BlMiJdLCJhdHRyaWJ1dGVzIjp7ImF0dHIxIjoidmFsdWUxIiwiYXR0cjIiOiIifX0.K3kRT2mbAtR_I7A_7zpRHnW_KggBHbCEqjYjhAt6mzqWcYQuEhVfbQerErqD3coyxoC4zIKsdnubtgl6ijWTsw"
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
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token.Username()).To(Equal("test-user"))
			})

			It("sets attributes", func() {
				token, _ := subject.ParseRequestAuth(req)

				Expect(token.HasAttribute("attr1")).To(BeTrue())
				Expect(token.GetAttribute("attr1")).To(Equal("value1"))
				Expect(token.HasAttribute("attr2")).To(BeTrue())
				Expect(token.GetAttribute("attr2")).To(Equal(""))
			})

			It("sets scopes", func() {
				token, _ := subject.ParseRequestAuth(req)

				Expect(token.HasAttribute("scope1")).To(BeTrue())
				Expect(token.GetAttribute("scope1")).To(BeTrue())
				Expect(token.HasAttribute("scope2")).To(BeTrue())
				Expect(token.GetAttribute("scope2")).To(BeTrue())
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

					token, err := subject.ParseRequestAuth(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid signing method"))
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

					token, err := subject.ParseRequestAuth(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Parsing claims"))
					Expect(token).To(BeNil())
				})
			})
		})

		Context("missing Authorization header", func() {
			It("bypasses", func() {
				req := http.Request{}

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

				token, err := subject.ParseRequestAuth(req)

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})
		})
	})
})
