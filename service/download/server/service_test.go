package server_test

import (
	"net/http"

	"github.com/dpb587/ssoca/server/service"
	. "github.com/dpb587/ssoca/service/download/server"
	svcconfig "github.com/dpb587/ssoca/service/download/server/config"

	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

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

	Context("basics", func() {
		var fs boshsysfakes.FakeFileSystem

		BeforeEach(func() {
			fs = *boshsysfakes.NewFakeFileSystem()

			subject = NewService("fake-name", svcconfig.Config{}, &fs)
		})

		Describe("Name", func() {
			It("works", func() {
				Expect(subject.Name()).To(Equal("fake-name"))
			})
		})

		Describe("Metadata", func() {
			It("works", func() {
				Expect(subject.Metadata()).To(BeNil())
			})
		})

		Describe("VerifyAuthorization", func() {
			It("works", func() {
				err := subject.VerifyAuthorization(http.Request{}, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
