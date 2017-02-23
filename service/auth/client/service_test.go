package client_test

import (
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
			Expect(subject.Type()).To(Equal("auth"))
		})

		It("Version", func() {
			Expect(subject.Version()).ToNot(Equal(""))
		})

		It("Description", func() {
			Expect(subject.Description()).ToNot(Equal(""))
		})
	})
})
