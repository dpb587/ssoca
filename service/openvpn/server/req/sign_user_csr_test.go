package req_test

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	"github.com/dpb587/ssoca/certauth/memory/memoryfakes"
	"github.com/dpb587/ssoca/server"
	"github.com/dpb587/ssoca/service/openvpn/api"
	. "github.com/dpb587/ssoca/service/openvpn/server/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SignUserCSR", func() {
	var subject SignUserCSR

	// certstrap request-cert --cn test --key-bits 1024 --passphrase ''
	var usr1csrStr = `-----BEGIN CERTIFICATE REQUEST-----
MIIBTjCBuAIBADAPMQ0wCwYDVQQDEwR0ZXN0MIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQCaQgeKTpVjtFwd5fiH1bvzovE3KrFiT/slPgWFYFlJT9TBQ20nwicw
qv+Lfbnr1vKjUk7aFkU3ihB+qcYYk/J0kuXcZFXz53jMD9kR7wDCNbpmvzyoejIo
0NHelzDHDyl9zq/xn6GEDDWIx5kIWzx/rEri4uve+mC8/uS47Wt+yQIDAQABoAAw
DQYJKoZIhvcNAQELBQADgYEAgVgNTHiD0yihWzVy59X1tMNVc+KBLLjzZZXRR4mk
L2e2xkgR/FAcp3ndKzk4tfak94VohbGvXzxieTtvDpfMUEYpWf7FQzPUBaZuebkC
aLUVg3Hw2wG7zZry4BtFfQnl8RDqqEsnj+41PUX2/eDbxd3pDr/61rUWqfQir1Xt
vqQ=
-----END CERTIFICATE REQUEST-----`

	pemToCertificate := func(bytes []byte) x509.Certificate {
		pem, _ := pem.Decode(bytes)
		if pem == nil {
			panic("Failed decoding PEM")
		}

		certificate, err := x509.ParseCertificate(pem.Bytes)
		if err != nil {
			panic(err)
		}

		return *certificate
	}

	Describe("Route", func() {
		It("returns", func() {
			Expect(subject.Route()).To(Equal("sign-user-csr"))
		})
	})

	Describe("Execute", func() {
		var fakecertauth certauthfakes.FakeProvider
		var realcertauth certauth.Provider
		var token auth.Token
		var loggerContext logrus.Fields

		BeforeEach(func() {
			loggerContext = logrus.Fields{
				"custom": "fake",
			}

			token = auth.Token{ID: "fake-user"}
			fakecertauth = certauthfakes.FakeProvider{}
			realcertauth = memoryfakes.CreateMock1()

			subject = SignUserCSR{
				Validity:    time.Duration(3600 * time.Second),
				CertAuth:    &fakecertauth,
				BaseProfile: "fake-profile",
			}
		})

		It("works", func() {
			fakecertauth.SignCertificateStub = realcertauth.SignCertificate

			res, err := subject.Execute(
				&token,
				api.SignUserCSRRequest{
					CSR: usr1csrStr,
				},
				loggerContext,
			)

			Expect(err).ToNot(HaveOccurred())

			Expect(res.Certificate).ToNot(Equal(""))
			pemToCertificate([]byte(res.Certificate))

			Expect(res.Profile).To(ContainSubstring("\nremap-usr1 SIGTERM\n"))
			Expect(res.Profile).To(ContainSubstring(fmt.Sprintf("<cert>\n%s\n</cert>", res.Certificate)))

			cert, _, innerLoggerContext := fakecertauth.SignCertificateArgsForCall(0)
			Expect(cert.SerialNumber).ToNot(Equal(0))
			Expect(cert.Subject.Organization).To(HaveLen(1))
			Expect(cert.Subject.Organization).To(ContainElement(Equal("ssoca/0.1.0")))
			Expect(cert.Subject.CommonName).To(Equal("fake-user"))
			Expect(cert.NotBefore.Unix()).To(BeNumerically(">", time.Now().Add(-10*time.Second).Unix()))
			Expect(cert.NotBefore.Unix()).To(BeNumerically("<", time.Now().Unix()))
			Expect(cert.NotAfter.Unix()).To(BeNumerically(">", time.Now().Add(3590*time.Second).Unix()))
			Expect(cert.NotAfter.Unix()).To(BeNumerically("<", time.Now().Add(3610*time.Second).Unix()))
			Expect(cert.KeyUsage & x509.KeyUsageKeyEncipherment).To(Equal(x509.KeyUsageKeyEncipherment))
			Expect(cert.KeyUsage & x509.KeyUsageDigitalSignature).To(Equal(x509.KeyUsageDigitalSignature))
			Expect(cert.ExtKeyUsage).To(HaveLen(1))
			Expect(cert.ExtKeyUsage).To(ContainElement(x509.ExtKeyUsageClientAuth))
			Expect(cert.BasicConstraintsValid).To(Equal(true))
			Expect(innerLoggerContext["custom"]).To(Equal("fake"))
		})

		Context("invalid csr", func() {
			Context("invalid format", func() {
				It("errors", func() {
					_, err := subject.Execute(
						&token,
						api.SignUserCSRRequest{
							CSR: "invalid",
						},
						loggerContext,
					)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Decoding CSR"))

					apiError, ok := err.(server.APIError)
					Expect(ok).To(BeTrue())

					Expect(apiError.Status).To(Equal(400))
					Expect(apiError.PublicError).To(Equal("Failed to decode certificate signing request"))
				})
			})

			Context("invalid data", func() {
				It("errors", func() {
					_, err := subject.Execute(
						&token,
						api.SignUserCSRRequest{
							CSR: `-----BEGIN CERTIFICATE REQUEST-----
MIIBTjCBuAIBADAPMQ0wCwYDVQQDEwR0ZXN0MIGfMA0GCSqGSIb3DQEBAQUAA4GN
-----END CERTIFICATE REQUEST-----`,
						},
						loggerContext,
					)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Parsing CSR"))

					apiError, ok := err.(server.APIError)
					Expect(ok).To(BeTrue())

					Expect(apiError.Status).To(Equal(400))
					Expect(apiError.PublicError).To(Equal("Failed to parse certificate signing request"))
				})
			})
		})

		Context("certauth failure", func() {
			Context("signing csr", func() {
				It("errors", func() {
					fakecertauth.SignCertificateReturns([]byte{}, errors.New("fake-err"))

					_, err := subject.Execute(
						&token,
						api.SignUserCSRRequest{
							CSR: usr1csrStr,
						},
						loggerContext,
					)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-err"))
					Expect(err.Error()).To(ContainSubstring("Signing certificate"))
				})
			})

			Context("getting ca certificate", func() {
				It("errors", func() {
					fakecertauth.SignCertificateStub = realcertauth.SignCertificate
					fakecertauth.GetCertificatePEMReturns("", errors.New("fake-err"))

					_, err := subject.Execute(
						&token,
						api.SignUserCSRRequest{
							CSR: usr1csrStr,
						},
						loggerContext,
					)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-err"))
					Expect(err.Error()).To(ContainSubstring("Loading CA certificate"))
				})
			})
		})
	})
})
