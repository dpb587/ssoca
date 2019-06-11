package service_test

import (
	. "github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/client/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultManager", func() {
	var subject Manager

	BeforeEach(func() {
		subject = NewDefaultManager()
	})

	Describe("Get", func() {
		It("retrieves services", func() {
			service := &servicefakes.FakeService{}
			service.TypeReturns("test1")
			service.NameReturns("fake-name")

			subject.Add(service)

			get, err := subject.Get("test1", "fake-name")

			Expect(err).ToNot(HaveOccurred())
			Expect(get).To(Equal(service))
		})

		Context("non-existant service", func() {
			It("errors", func() {
				_, err := subject.Get("test1", "nop")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unrecognized type: test1"))
			})
		})
	})
})
