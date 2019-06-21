package server_test

import (
	"github.com/dpb587/ssoca/service/httpauth"
	. "github.com/dpb587/ssoca/service/httpauth/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var subject ServiceFactory

	BeforeEach(func() {
		subject = NewServiceFactory()
	})

	Describe("Type", func() {
		It("works", func() {
			Expect(subject.Type()).To(Equal(httpauth.Type))
		})
	})

	Describe("Create", func() {
		It("works", func() {
			service, err := subject.Create("name1", map[string]interface{}{
				"invalid": true,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(service.Name()).To(Equal("name1"))
		})
	})
})
