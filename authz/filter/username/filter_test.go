package username_test

import (
	. "github.com/dpb587/ssoca/authz/filter/username"

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
					"is": "dberger",
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(req).ToNot(BeNil())

				treq, ok := req.(Requirement)
				Expect(ok).To(BeTrue())
				Expect(treq.Is).To(Equal("dberger"))
			})
		})

		Context("with invalid config", func() {
			Context("with missing property (is)", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(": is"))
				})
			})

			Context("with incorrect types", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{
						"is": []string{"dberger"},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Loading config"))
				})
			})
		})
	})
})
