package githubauth_test

import (
	. "github.com/dpb587/ssoca/service/githubauth"

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
			Expect(svc.Type()).To(Equal("github_authn"))
		})
	})

	Describe("Version", func() {
		It("works", func() {
			Expect(svc.Version()).ToNot(Equal(""))
		})
	})
})
