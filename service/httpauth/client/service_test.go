package client_test

import (
	"github.com/dpb587/ssoca/service/httpauth"
	. "github.com/dpb587/ssoca/service/httpauth/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("basics", func() {
		var subject Service

		BeforeEach(func() {
			subject = Service{}
		})

		It("Type", func() {
			Expect(subject.Type()).To(Equal(httpauth.Type))
		})

		It("Version", func() {
			Expect(subject.Version()).ToNot(Equal(""))
		})
	})
})
