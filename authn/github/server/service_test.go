package server_test

import (
	"net/http"

	svcconfig "github.com/dpb587/ssoca/authn/github/config"
	. "github.com/dpb587/ssoca/authn/github/server"
	oauth2support "github.com/dpb587/ssoca/authn/support/oauth2"
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
			subject = NewService("test1", svcconfig.Config{}, oauth2support.Backend{})
		})

		Describe("Name", func() {
			It("works", func() {
				Expect(subject.Name()).To(Equal("test1"))
			})
		})

		Describe("Metadata", func() {
			It("works", func() {
				Expect(subject.Metadata()).To(BeNil())
			})
		})

		Describe("IsAuthorized", func() {
			It("works", func() {
				authz, err := subject.IsAuthorized(http.Request{}, nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(authz).To(BeTrue())
			})
		})
	})
})
