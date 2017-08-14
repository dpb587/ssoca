package config_test

import (
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	. "github.com/dpb587/ssoca/service/openvpn/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var subject Config

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
