package uaa_test

import (
	. "github.com/dpb587/ssoca/authn/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var svc Service

	BeforeEach(func() {
		svc = Service{}
	})

	Describe("Type", func() {
		It("works", func() {
			Expect(svc.Type()).To(Equal("uaa_authn"))
		})
	})

	Describe("Version", func() {
		It("works", func() {
			Expect(svc.Version()).ToNot(Equal(""))
		})
	})
})
