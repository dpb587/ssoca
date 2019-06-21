package client_test

import (
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/service/uaaauth"
	. "github.com/dpb587/ssoca/service/uaaauth/client"

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
			Expect(subject.Type()).To(Equal(uaaauth.Type))
		})

		It("Version", func() {
			Expect(subject.Version()).ToNot(Equal(""))
		})
	})
})
