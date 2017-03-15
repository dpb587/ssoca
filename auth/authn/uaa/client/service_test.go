package client_test

import (
	. "github.com/dpb587/ssoca/auth/authn/uaa/client"

	"github.com/dpb587/ssoca/client/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Describe("interface", func() {
		It("github.com/dpb587/ssoca/client/service.AuthService", func() {
			var _ service.AuthService = (*Service)(nil)
		})
	})

	Context("basics", func() {
		var subject Service

		BeforeEach(func() {
			subject = Service{}
		})

		It("Type", func() {
			Expect(subject.Type()).To(Equal("uaa_authn"))
		})

		It("Version", func() {
			Expect(subject.Version()).ToNot(Equal(""))
		})

		It("Description", func() {
			Expect(subject.Description()).ToNot(Equal(""))
		})
	})
})
