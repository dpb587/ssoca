package server_test

import (
	"net/http"

	. "github.com/dpb587/ssoca/authn/http/server"

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

		service = NewService(
			"auth",
			Config{
				Users: []UserConfig{
					{
						Username: "user1",
						Password: "pass1",
						Attributes: map[string]interface{}{
							"scope1": true,
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
			Expect(token.Username()).To(Equal("user1"))
			Expect(token.Attributes()).To(HaveLen(1))
			Expect(token.Attributes()["scope1"]).To(BeTrue())
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
