package server_test

import (
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	"github.com/dpb587/ssoca/service/openvpn"
	. "github.com/dpb587/ssoca/service/openvpn/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceFactory", func() {
	var subject ServiceFactory

	Describe("Type", func() {
		It("returns", func() {
			Expect(subject.Type()).To(Equal(openvpn.Type))
		})
	})

	Describe("Create", func() {
		var caManager certauth.Manager

		BeforeEach(func() {
			caManager = certauth.NewDefaultManager()

			certauth := certauthfakes.FakeProvider{}
			certauth.NameReturns("default")
			caManager.Add(&certauth)

			subject = NewServiceFactory(caManager)
		})

		It("remarshals configuration", func() {
			provider, err := subject.Create("name1", map[string]interface{}{
				"profile": "something",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(provider.Name()).To(Equal("name1"))
		})

		Context("invalid certauth", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certauth": "unknown",
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("getting certificate authority"))
			})
		})

		Context("invalid validity duration", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"validity": "525,600 minutes",
				})

				Expect(err).To(HaveOccurred())
			})
		})

		Context("invalid yaml", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"profile": map[string]interface{}{},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("loading config"))
			})
		})
	})
})
