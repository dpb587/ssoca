package scope_test

import (
	. "github.com/dpb587/ssoca/auth/authz/filter/scope"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("Create", func() {
		var subject Filter

		BeforeEach(func() {
			subject = Filter{}
		})

		Context("with valid config", func() {
			It("creates", func() {
				req, err := subject.Create(map[string]interface{}{
					"present": "scope1",
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(req).ToNot(BeNil())

				treq, ok := req.(Requirement)
				Expect(ok).To(BeTrue())
				Expect(treq.Present).To(Equal("scope1"))
			})
		})

		Context("with invalid config", func() {
			Context("with missing property (present)", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(": present"))
				})
			})

			Context("with incorrect types", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{
						"present": []string{"else"},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("loading config"))
				})
			})
		})
	})
})
