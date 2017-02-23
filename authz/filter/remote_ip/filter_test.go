package remote_ip_test

import (
	. "github.com/dpb587/ssoca/authz/filter/remote_ip"

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
					"within": "192.0.2.0/24",
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(req).ToNot(BeNil())

				treq, ok := req.(Requirement)
				Expect(ok).To(BeTrue())
				Expect(treq.WithinRaw).To(Equal("192.0.2.0/24"))
				Expect(treq.Within.String()).To(Equal("192.0.2.0/24"))
			})

			Context("with IP format", func() {
				Context("IPv4", func() {
					It("appends /32", func() {
						req, err := subject.Create(map[string]interface{}{
							"within": "192.0.2.127",
						})

						Expect(err).ToNot(HaveOccurred())
						Expect(req).ToNot(BeNil())

						treq, ok := req.(Requirement)
						Expect(ok).To(BeTrue())
						Expect(treq.WithinRaw).To(Equal("192.0.2.127"))
						Expect(treq.Within.String()).To(Equal("192.0.2.127/32"))
					})
				})

				Context("IPv6", func() {
					It("appends /128", func() {
						req, err := subject.Create(map[string]interface{}{
							"within": "abcd:ef01:2345:6789:abcd:ef01:2345:6789",
						})

						Expect(err).ToNot(HaveOccurred())
						Expect(req).ToNot(BeNil())

						treq, ok := req.(Requirement)
						Expect(ok).To(BeTrue())
						Expect(treq.WithinRaw).To(Equal("abcd:ef01:2345:6789:abcd:ef01:2345:6789"))
						Expect(treq.Within.String()).To(Equal("abcd:ef01:2345:6789:abcd:ef01:2345:6789/128"))
					})
				})
			})
		})

		Context("with invalid config", func() {
			Context("with missing property (within)", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(": within"))
				})
			})

			Context("with incorrect types", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{
						"within": map[string]interface{}{},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Loading config"))
				})
			})

			Context("with invalid format", func() {
				It("errors", func() {
					_, err := subject.Create(map[string]interface{}{
						"within": "something",
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Parsing CIDR"))
				})
			})
		})
	})
})
