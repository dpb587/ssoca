package server_test

import (
	"net/http"

	svcapi "github.com/dpb587/ssoca/auth/authn/uaa/api"
	svcconfig "github.com/dpb587/ssoca/auth/authn/uaa/config"
	. "github.com/dpb587/ssoca/auth/authn/uaa/server"

	"github.com/dpb587/ssoca/server/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var subject Service

	Describe("interface", func() {
		It("github.com/dpb587/ssoca/server/service.AuthService", func() {
			var _ service.AuthService = (*Service)(nil)
		})
	})

	Context("Basics", func() {
		BeforeEach(func() {
			subject = NewService("test1", svcconfig.Config{
				URL:           "my-url",
				CACertificate: "my-certificate",
				PublicKey:     "public-key",
				ClientID:      "client-id",
				ClientSecret:  "client-secret",
				Prompts: []string{
					"one",
					"two",
				},
			})
		})

		Describe("Name", func() {
			It("works", func() {
				Expect(subject.Name()).To(Equal("test1"))
			})
		})

		Describe("Metadata", func() {
			It("works", func() {
				Expect(subject.Metadata()).To(Equal(svcapi.Metadata{
					URL:           "my-url",
					CACertificate: "my-certificate",
					ClientID:      "client-id",
					ClientSecret:  "client-secret",
					Prompts: []string{
						"one",
						"two",
					},
				}))
			})
		})

		Describe("GetRoutes", func() {
			It("works", func() {
				Expect(subject.GetRoutes()).To(BeNil())
			})
		})

		Describe("VerifyAuthorization", func() {
			It("works", func() {
				err := subject.VerifyAuthorization(http.Request{}, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
