package config_test

import (
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	. "github.com/dpb587/ssoca/service/ssh/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var subject Config

	// Describe("ApplyDefaults", func() {
	// 	Describe("with defaults", func() {
	// 		BeforeEach(func() {
	// 			subject.ApplyDefaults()
	// 		})
	//
	// 		It("defaults extensions", func() {
	// 			Expect(subject.Extensions).To(Equal(ExtensionDefaults))
	// 		})
	// 	})
	//
	// 	Describe("not really dealing with defaults but its convenient so...", func() {
	// 		It("ignores extensions if ssoca-no-defaults is used", func() {
	// 			subject.Extensions = Extensions{ExtensionNoDefaults}
	// 			subject.ApplyDefaults()
	//
	// 			Expect(subject.Extensions).To(HaveLen(0))
	// 		})
	// 	})
	// })

	Context("when certauth is not configured", func() {
		It("defaults to default", func() {
			fakecertauth := &certauthfakes.FakeProvider{}
			fakecertauth.GetCertificatePEMReturns("fake-certificate-pem", nil)

			fakecertauths := &certauthfakes.FakeManager{}
			fakecertauths.GetReturns(fakecertauth, nil)

			subject.CertAuth = certauth.NewConfigValue(fakecertauths)
			subject.ApplyDefaults()

			pem, err := subject.CertAuth.GetCertificatePEM()
			Expect(err).ToNot(HaveOccurred())
			Expect(pem).To(Equal("fake-certificate-pem"))
		})
	})
})
