package server_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	svcconfig "github.com/dpb587/ssoca/auth/authn/http/config"
	. "github.com/dpb587/ssoca/auth/authn/http/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var service Service
	var request http.Request

	BeforeEach(func() {
		request = http.Request{
			Header: http.Header{},
		}
		name := "test 1"

		service = NewService(
			"auth",
			svcconfig.Config{
				Users: []svcconfig.User{
					{
						Username: "user1",
						Password: "pass1",
						Groups: []string{
							"scope1",
						},
						Attributes: map[auth.TokenAttribute]*string{
							auth.TokenNameAttribute: &name,
						},
					},
				},
			},
		)
	})

	Describe("ParseRequestAuth", func() {
		It("creates tokens", func() {
			request.SetBasicAuth("user1", "pass1")

			token, err := service.ParseRequestAuth(request)

			Expect(err).ToNot(HaveOccurred())

			Expect(token).ToNot(BeNil())
			Expect(token.ID).To(Equal("user1"))
			Expect(token.Username()).To(Equal("user1"))
			Expect(token.Name()).To(Equal("test 1"))
			Expect(token.Groups).To(HaveLen(1))
			Expect(token.Groups).To(ContainElement("scope1"))
		})

		Context("without authentication", func() {
			It("has no token or error", func() {
				token, err := service.ParseRequestAuth(request)

				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(BeNil())
			})
		})

		Context("with invalid user/pass", func() {
			It("errors", func() {
				request.SetBasicAuth("baduser", "badpass")

				token, err := service.ParseRequestAuth(request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).ToNot(ContainSubstring("badpass"))

				Expect(token).To(BeNil())
			})
		})

		Context("with invalid password", func() {
			It("errors", func() {
				request.SetBasicAuth("user1", "badpass")

				token, err := service.ParseRequestAuth(request)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).ToNot(ContainSubstring("badpass"))

				Expect(token).To(BeNil())
			})
		})
	})
})
