package server_test

import (
	"net/http"

	"github.com/dpb587/ssoca/server/service"
	svcconfig "github.com/dpb587/ssoca/service/openvpn/config"
	. "github.com/dpb587/ssoca/service/openvpn/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var subject Service

	Describe("interface", func() {
		It("github.com/dpb587/ssoca/server/service.Service", func() {
			var _ service.Service = (*Service)(nil)
		})
	})

	Context("Basics", func() {
		BeforeEach(func() {
			subject = NewService("test1", svcconfig.Config{})
		})

		Describe("Name", func() {
			It("works", func() {
				Expect(subject.Name()).To(Equal("test1"))
			})
		})

		Describe("Metadata", func() {
			It("works", func() {
				Expect(subject.Metadata()).To(BeNil())
			})
		})

		Describe("IsAuthorized", func() {
			It("works", func() {
				authz, err := subject.IsAuthorized(http.Request{}, nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(authz).To(BeTrue())
			})
		})
	})
})
