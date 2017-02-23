package service_test

import (
	"errors"

	. "github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultFactory", func() {
	Describe("Create", func() {
		It("works", func() {
			service := &servicefakes.FakeService{}
			serviceOptions := map[string]interface{}{"key1": "val1"}

			serviceFactory := &servicefakes.FakeServiceFactory{}
			serviceFactory.TypeReturns("fake1")
			serviceFactory.CreateReturns(service, nil)

			factory := NewDefaultFactory()
			factory.Register(serviceFactory)

			svc, err := factory.Create("fake1", "test1", serviceOptions)

			Expect(err).ToNot(HaveOccurred())
			Expect(svc).To(Equal(service))

			arg0, arg1 := serviceFactory.CreateArgsForCall(0)
			Expect(arg0).To(Equal("test1"))
			Expect(arg1).To(Equal(serviceOptions))
		})

		Context("unknown service", func() {
			It("errors", func() {
				factory := NewDefaultFactory()

				_, err := factory.Create("unknown1", "test1", map[string]interface{}{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown1"))
			})
		})

		Context("factory errors", func() {
			It("errors", func() {
				serviceFactory := &servicefakes.FakeServiceFactory{}
				serviceFactory.TypeReturns("fake1")
				serviceFactory.CreateReturns(nil, errors.New("fake-error1"))

				factory := NewDefaultFactory()
				factory.Register(serviceFactory)

				_, err := factory.Create("fake1", "test1", map[string]interface{}{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Creating service fake1[test1]"))
			})
		})
	})
})
