package server_test

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	internaltests "github.com/dpb587/ssoca/service/uaaauth/internal/tests"
	. "github.com/dpb587/ssoca/service/uaaauth/server"
	svcconfig "github.com/dpb587/ssoca/service/uaaauth/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var subject *Service

	BeforeEach(func() {
		subject = NewService(
			"auth",
			svcconfig.Config{
				PublicKey: internaltests.SharedPublicKey,
			},
		)
	})

	Describe("ParseRequestAuth", func() {
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
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token.ID).To(Equal("fake-user1"))
			})

			It("sets attributes", func() {
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(*token.Attributes[auth.TokenAttribute("username")]).To(Equal("fake-user1"))
			})

			It("sets groups", func() {
				token, err := subject.ParseRequestAuth(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(token.Groups).To(ContainElement("scope1"))
				Expect(token.Groups).To(ContainElement("scope2"))
			})
		})

		Context("invalid token data", func() {
			Context("with algorithm none", func() {
				var tokenString = "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJleHAiOjE4MjkwMDE2MDAsImp0aSI6ImZha2UtdXVpZCIsImlhdCI6MTU0NTAwNDgwMCwiaXNzIjoiZmFrZS1vcmlnaW4iLCJuYmYiOjE1NDUwMDQ4MDAsImF1ZCI6WyJmYWtlLWF1ZGllbmNlMSJdLCJ1c2VyX25hbWUiOiJmYWtlLXVzZXIxIiwic2NvcGUiOlsic2NvcGUxIiwic2NvcGUyIl19."

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
					Expect(err.Error()).To(ContainSubstring("no signing method used"))

					errapi, errok := err.(apierr.Error)
					Expect(errok).To(BeTrue())
					Expect(errapi.Status).To(Equal(http.StatusForbidden))

					Expect(token).To(BeNil())
				})
			})

			Context("with invalid signature", func() {
				var tokenString = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MjkwMDE2MDAsImp0aSI6ImZha2UtdXVpZCIsImlhdCI6MTU0NTAwNDgwMCwiaXNzIjoiZmFrZS1vcmlnaW4iLCJuYmYiOjE1NDUwMDQ4MDAsImF1ZCI6WyJmYWtlLWF1ZGllbmNlMSJdLCJ1c2VyX25hbWUiOiJmYWtlLXVzZXIxIiwic2NvcGUiOlsic2NvcGUxIiwic2NvcGUyIl19.bTDTsI8zxPQ09XsY07ZqG2uqoSVHhSZo4q_XgPbRcHLkNVlGU6AIJXgteCBLThFkAQfDmxa6chCl2AMwT4w8QhfT4j0dFk1SvdWdlMwA-RpoW7ge-ei2BCUk9dCd0gos4XoqC3XNREWEui7qLdawgduPcfM67fx9Q0wPLCQflFjp7jzoL4-0LhGTRlxtC9cL06ybq9IH6JaqQGtOyj-Q4wHlZK8LNTXnRE8BfOJlX9fr0UzRoZE4vLTFpZscAxrEcJ4nuBPcWAEmTY-X1fivoPNT8gvCg8I_xph_RUsiqHaRDko6pSWEe5gJ41xdjpyXY_xW5Ipa-UpEZxTPMc9GVA"

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
					Expect(err.Error()).To(ContainSubstring("parsing claims"))

					errapi, errok := err.(apierr.Error)
					Expect(errok).To(BeTrue())
					Expect(errapi.Status).To(Equal(http.StatusForbidden))

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
