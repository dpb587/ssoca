package public_test

import (
	. "github.com/dpb587/ssoca/auth/authz/filter/public"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("Create", func() {
		var subject Filter

		BeforeEach(func() {
			subject = Filter{}
		})

		It("creates", func() {
			req, err := subject.Create(map[string]interface{}{})

			Expect(err).ToNot(HaveOccurred())
			Expect(req).ToNot(BeNil())

			_, ok := req.(Requirement)
			Expect(ok).To(BeTrue())
		})
	})
})
