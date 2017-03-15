package filter_test

import (
	. "github.com/dpb587/ssoca/auth/authz/filter"

	"github.com/dpb587/ssoca/auth/authz/filter/filterfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultManager", func() {
	Describe("Get", func() {
		It("returns valid filters", func() {
			filter1 := &filterfakes.FakeFilter{}

			manager := NewDefaultManager()
			manager.Add("fake1", filter1)

			f, err := manager.Get("fake1")

			Expect(err).ToNot(HaveOccurred())
			Expect(f).To(Equal(filter1))
		})

		Context("unknown filters", func() {
			It("errors", func() {
				manager := NewDefaultManager()

				f, err := manager.Get("unknown")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unrecognized filter"))
				Expect(f).To(BeNil())
			})
		})
	})

	Describe("Filters", func() {
		It("returns filters", func() {
			filter1 := &filterfakes.FakeFilter{}

			manager := NewDefaultManager()
			manager.Add("fake1", filter1)

			Expect(manager.Filters()).To(Equal([]string{"fake1"}))
		})
	})
})
