package management_test

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	. "github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"
	"github.com/dpb587/ssoca/service/openvpn/client/profile/profilefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultHandler", func() {
	var subject *DefaultHandler
	var fakeProfileManager *profilefakes.FakeManager
	var fakeWriter *bytes.Buffer

	var test1key *rsa.PrivateKey
	var test1keyStr = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCqAzEMN8rybTZMLfUjnrcXCPTAY7uYQHA1qRAcO02jJjr0NuxY
2eleYf31uRnyfXJsWAiecaZlwp52qttwCTtharJMgcs9Lr5Z07lUUUdOy93CHx8y
dlgKJAHCRWUtIXEAq0F2zm4Nlr98cGgaARMwvXTeRfXpkEQmeArAI4ntpwIDAQAB
AoGADWa/AQWM29s8AnlE74/dQtWT5W53JSM/NRukh3UtQ4UJ9KI3szFKMgRrbmku
4Gx/DodJ9qNiyHa04wnIzmYL5hr6OmGUUHDnBK8ZtLxzlHfthcOYJONPGOGgBdwG
zWRxFzNwnpqWyAS1G2yJln6wlN04grxAm3GnKTOMEYW8hUECQQDfFJALHF6aYsk+
6E0649bjBuchVy+pFFKamCn2/ZTqzULXFFACRSD4MSS/FjwCQkNeWC+4R2GrY8eJ
axqIL563AkEAwxnaZ3wf1RpfBUla08VxMcjMEc5UfsU+Y5tfnPc7rryJX6Hmg73B
uvHXj8VHVcfuLJjeSVQocEsrKW6I3+84kQJBAKYwFGsilFuhUlk6CCbiC3kP8GoH
IKtuR2eCCmlFWoZdqfi+2igGxdwACGcOsl/ga33CZrJ7AwkCiWkXUCm6iBsCQGQJ
qZdOafQXJYnMZyoXH0drslee+GxYLvlb/da6XnvmaHoExfHfJqr4vpMVkNJHRbTQ
XYo0ANgzcto3ty87tkECQQDN46eAjb9xSXFYLO/ILlpr3QU71v8l1zheGkuBNYOu
ZNYBM+NpfAXTMgHSuWnIkZSoSoV4ZcYTkJ6zslGbkVyH
-----END RSA PRIVATE KEY-----`

	BeforeEach(func() {
		fakeProfileManager = &profilefakes.FakeManager{}
		fakeWriter = bytes.NewBuffer(nil)

		test1keyPEM, _ := pem.Decode([]byte(test1keyStr))
		if test1keyPEM == nil {
			panic(errors.New("Failed decoding private key PEM"))
		}

		var err error

		test1key, err = x509.ParsePKCS1PrivateKey(test1keyPEM.Bytes)
		if err != nil {
			panic(err)
		}

		subject = NewDefaultHandler(fakeProfileManager)
	})

	BehavesLikeSimpleHandler := func(cb ServerHandlerCallback) {
		rcb, err := cb(nil, "SUCCESS: data")
		Expect(err).ToNot(HaveOccurred())
		Expect(rcb).To(BeNil())

		rcb, err = cb(nil, "FAILURE: data")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Bad management command result"))
	}

	Describe("NeedCertificate", func() {
		It("propagates errors", func() {
			fakeProfileManager.GetProfileReturns(profile.Profile{}, errors.New("fake-err1"))

			_, err := subject.NeedCertificate(fakeWriter, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
			Expect(err.Error()).To(ContainSubstring("Retrieving profile"))
		})

		It("responds with certificate", func() {
			fakeProfileManager.GetProfileReturns(profile.NewProfile("base\nconfig", test1key, []byte("fake\ncertificate")), nil)

			cb, err := subject.NeedCertificate(fakeWriter, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeWriter.String()).To(Equal("certificate\nfake\ncertificate\nEND\n"))

			BehavesLikeSimpleHandler(cb)
		})

		Context("openvpn connection reattempts", func() {
			It("renews certificates after 3 rapid failures", func() {
				fakeProfileManager.GetProfileReturns(profile.NewProfile("base\nconfig", test1key, []byte("fake\ncertificate")), nil)
				fakeProfileManager.RenewReturns(nil)

				for i := 0; i < 2; i++ {
					_, err := subject.NeedCertificate(fakeWriter, "")
					Expect(err).ToNot(HaveOccurred())
				}

				fakeWriter.Reset()

				_, err := subject.NeedCertificate(fakeWriter, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeWriter.String()).To(Equal("certificate\nfake\ncertificate\nEND\n"))
				Expect(fakeProfileManager.RenewCallCount()).To(Equal(1))

				fakeWriter.Reset()

				// shouldn't renew a subsequent time
				_, err = subject.NeedCertificate(fakeWriter, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeWriter.String()).To(Equal("certificate\nfake\ncertificate\nEND\n"))
				Expect(fakeProfileManager.RenewCallCount()).To(Equal(1))
			})

			It("propagates renewal errors", func() {
				fakeProfileManager.GetProfileReturns(profile.NewProfile("base\nconfig", test1key, []byte("fake\ncertificate")), nil)
				fakeProfileManager.RenewReturns(errors.New("fake-err1"))

				for i := 0; i < 2; i++ {
					_, err := subject.NeedCertificate(fakeWriter, "")
					Expect(err).ToNot(HaveOccurred())
				}

				_, err := subject.NeedCertificate(fakeWriter, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err1"))
				Expect(err.Error()).To(ContainSubstring("Renewing profile"))
			})

			It("restarts openvpn after 5 rapid failures", func() {
				fakeProfileManager.GetProfileReturns(profile.NewProfile("base\nconfig", test1key, []byte("fake\ncertificate")), nil)
				fakeProfileManager.RenewReturns(nil)

				for i := 0; i < 4; i++ {
					_, err := subject.NeedCertificate(fakeWriter, "")
					Expect(fakeWriter.String()).To(Equal("certificate\nfake\ncertificate\nEND\n"))
					Expect(err).ToNot(HaveOccurred())

					fakeWriter.Reset()
				}

				_, err := subject.NeedCertificate(fakeWriter, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeWriter.String()).To(Equal("signal SIGHUP\n"))

				fakeWriter.Reset()

				// should continuously respond with SIGHUP
				_, err = subject.NeedCertificate(fakeWriter, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeWriter.String()).To(Equal("signal SIGHUP\n"))
			})
		})
	})

	Describe("SignRSA", func() {
		It("sends sighup on invalid certificate", func() {
			cb, err := subject.SignRSA(fakeWriter, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeWriter.String()).To(Equal("signal SIGHUP\n"))

			BehavesLikeSimpleHandler(cb)
		})

		It("errors on invalid base64 data", func() {
			fakeProfileManager.IsCertificateValidReturns(true)

			_, err := subject.SignRSA(fakeWriter, "fake$data")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("illegal base64 data at input byte 4"))
			Expect(err.Error()).To(ContainSubstring("Decoding signing token"))
		})

		It("errors when signing fails", func() {
			fakeProfileManager.IsCertificateValidReturns(true)
			fakeProfileManager.SignReturns(nil, errors.New("fake-err1"))

			_, err := subject.SignRSA(fakeWriter, "aW1hc3NvY2E=")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
			Expect(err.Error()).To(ContainSubstring("Signing token"))
		})

		It("signs data", func() {
			fakeProfileManager.IsCertificateValidReturns(true)
			fakeProfileManager.SignReturns([]byte("fake-signature"), nil)

			cb, err := subject.SignRSA(fakeWriter, "aW1hc3NvY2E=")
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeWriter.String()).To(Equal("rsa-sig\nZmFrZS1zaWduYXR1cmU=\nEND\n"))

			BehavesLikeSimpleHandler(cb)

			Expect(fakeProfileManager.SignCallCount()).To(Equal(1))
			Expect(fakeProfileManager.SignArgsForCall(0)).To(Equal([]byte("imassoca")))
		})
	})
})
