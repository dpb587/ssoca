package client_test

import (
	svc "github.com/dpb587/ssoca/service/auth"
	. "github.com/dpb587/ssoca/service/auth/client"

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
			Expect(subject.Type()).To(Equal(svc.Type))
		})

		It("Version", func() {
			Expect(subject.Version()).ToNot(Equal(""))
		})
	})
})
