package googleauth_test

import (
	. "github.com/dpb587/ssoca/service/googleauth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceType", func() {
	var svc ServiceType

	BeforeEach(func() {
		svc = ServiceType{}
	})

	Describe("Type", func() {
		It("works", func() {
			Expect(svc.Type()).To(Equal("google_authn"))
		})
	})

	Describe("Version", func() {
		It("works", func() {
			Expect(svc.Version()).ToNot(Equal(""))
		})
	})
})