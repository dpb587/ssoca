package or_test

import (
	"errors"

	"github.com/dpb587/ssoca/authz/filter"
	"github.com/dpb587/ssoca/authz/filter/filterfakes"
	. "github.com/dpb587/ssoca/authz/filter/or"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("Create", func() {
		var subject Filter
		var requirement1, requirement2 filterfakes.FakeRequirement

		BeforeEach(func() {
			manager := filter.NewDefaultManager()

			subject = NewFilter(&manager)

			requirement1 = filterfakes.FakeRequirement{}
			filterValid := filterfakes.FakeFilter{}
			filterValid.CreateReturns(&requirement1, nil)
			manager.Add("valid", &filterValid)

			requirement2 = filterfakes.FakeRequirement{}
			filterError := filterfakes.FakeFilter{}
			filterError.CreateReturns(&requirement2, errors.New("fake error"))
			manager.Add("error", &filterError)
		})

		Context("with valid config", func() {
			It("creates", func() {
				req, err := subject.Create([]filter.RequireConfig{
					filter.RequireConfig{
						"valid": map[string]interface{}{},
					},
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(req).ToNot(BeNil())

				treq, ok := req.(Requirement)
				Expect(ok).To(BeTrue())
				Expect(treq.Requirements).To(Equal([]filter.Requirement{&requirement1}))
			})
		})

		Context("with invalid config", func() {
			Context("bad formatting or structure", func() {
				Context("with non-config filter list", func() {
					It("errors", func() {
						_, err := subject.Create("something")

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Failed to parse 'or' config"))
					})
				})

				Context("with no keys in hash array item", func() {
					It("errors", func() {
						_, err := subject.Create([]filter.RequireConfig{
							filter.RequireConfig{},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("does not have 1 field"))
					})
				})

				Context("with multiple keys in hash array item", func() {
					It("errors", func() {
						_, err := subject.Create([]filter.RequireConfig{
							filter.RequireConfig{
								"one": map[string]interface{}{},
								"two": map[string]interface{}{},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("does not have 1 field"))
					})
				})
			})

			Context("invalid filter", func() {
				It("errors", func() {
					_, err := subject.Create([]filter.RequireConfig{
						filter.RequireConfig{
							"invalid": map[string]interface{}{},
						},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Loading filter 'invalid'"))
				})
			})

			Context("invalid filter configuration", func() {
				It("errors", func() {
					_, err := subject.Create([]filter.RequireConfig{
						filter.RequireConfig{
							"error": map[string]interface{}{},
						},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Creating requirement"))
				})
			})
		})
	})
})
