package service_test

import (
	. "github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultManager", func() {
	Describe("Get (and Add)", func() {
		It("works", func() {
			service := &servicefakes.FakeService{}
			service.NameReturns("fake1")

			manager := NewDefaultManager()
			manager.Add(service)

			get, err := manager.Get("fake1")

			Expect(err).ToNot(HaveOccurred())
			Expect(get).To(Equal(service))
		})

		Context("invalid service", func() {
			It("errors", func() {
				manager := NewDefaultManager()

				_, err := manager.Get("fake1")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(": fake1"))
			})
		})
	})

	Describe("Services", func() {
		It("works", func() {
			service1 := &servicefakes.FakeService{}
			service1.NameReturns("fake1")

			service2 := &servicefakes.FakeService{}
			service2.NameReturns("fake2")

			manager := NewDefaultManager()
			manager.Add(service1)
			manager.Add(service2)

			services := manager.Services()

			Expect(services).To(ContainElement("fake1"))
			Expect(services).To(ContainElement("fake2"))
		})
	})
})
