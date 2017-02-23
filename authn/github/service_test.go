package github_test

import (
	. "github.com/dpb587/ssoca/authn/github"

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
			Expect(svc.Type()).To(Equal("github_authn"))
		})
	})

	Describe("Version", func() {
		It("works", func() {
			Expect(svc.Version()).ToNot(Equal(""))
		})
	})
})
