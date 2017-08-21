package profile_test

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"

	. "github.com/dpb587/ssoca/service/openvpn/client/profile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dpb587/ssoca/service/openvpn/api"
	"github.com/dpb587/ssoca/service/openvpn/httpclient/httpclientfakes"
)

var _ = Describe("CreateManagerAndPrivateKey", func() {
	// lame
	It("executes", func() {
		CreateManagerAndPrivateKey(&httpclientfakes.FakeClient{}, "fake-service")
	})
})

var _ = Describe("DefaultManager", func() {
	var subject DefaultManager
	var fakeapiclient *httpclientfakes.FakeClient

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
	var test1crtValid = `-----BEGIN CERTIFICATE-----
MIICsDCCAhmgAwIBAgIJAJcAjLkz1BLWMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTcwODA1MDY1ODM3WhcNMjcwODAzMDY1ODM3WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDQ4i/Vwwvw8pmXOUx32QMhZGsEDZ+Ifk14yQVJ4ZaZG/Sc9qJsMWmETE75/Ego
9VvLmtnA2IzbIisC6WmUM5Nz7LfmU1ogpMSHLBOVxd63fThT2/SnPOeBQSDsKEyb
BnWWjFEii9esSfbzU8TU2Bmgl7CwCRw6cxZQyk4+aJ3iPwIDAQABo4GnMIGkMB0G
A1UdDgQWBBTG2fDW9rRiS3PgXNIMl/WTpRrU5jB1BgNVHSMEbjBsgBTG2fDW9rRi
S3PgXNIMl/WTpRrU5qFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUt
U3RhdGUxITAfBgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJAJcAjLkz
1BLWMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAGwKdS5ccuA/8PN+v
BzJE953p4RsFtblzF5kERrouKJSpXlX5xWeC8mfJdzEQPOpLyEQUbYe2JGyEEExP
Mu9vPAdIpIUf9rx3CA1N0idhhec28nvmBh2066y9H1FDjfSB75mkP2MPI+7PBYKw
hFOYEgUmyTKD8Cwy7VCz50M2kh4=
-----END CERTIFICATE-----` // 2027-08-03T06:58:37Z
	var test1crtExpired = `-----BEGIN CERTIFICATE-----
MIICsDCCAhmgAwIBAgIJAIUnEIDGekg3MA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTcwODA1MDY1NTI1WhcNMTcwODA0MDY1NTI1WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDUSHqB5mm3XvwJtiIEeoM4eWzMmMQTPUgPYvkk4uMJR/b8VsVGwEX0DwY5sO59
1l4Kovg5cupyE52WntAHNfFjbAHfCfuRNhRTUVGkzdrNhe0QxPVXdjG3RRgzM32C
07r3Vwz2+MZEDehrvOL2dS/MN2pKEG7SGFh16mM3Fla8aQIDAQABo4GnMIGkMB0G
A1UdDgQWBBT727RDhGxKKc3cO12BaNfC/WixiDB1BgNVHSMEbjBsgBT727RDhGxK
Kc3cO12BaNfC/WixiKFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUt
U3RhdGUxITAfBgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJAIUnEIDG
ekg3MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADgYEAYGeTyXozAo9hhmAa
/2hQqQinuN3R1kx53P4+EyzTQMX3bPxp6ppPEnQPg1d38kKRKxUl/3WHIOndfuVH
8bjlsc7bv9VmtUy3bYW/jnJTU9kLwLUjNDc34QpBdW1ofRSaLY+p2X0TlI5RTSMW
D5ASZtCvQ6NPMPmlb2p2qCB8Ljk=
-----END CERTIFICATE-----` // 2017-08-04T06:55:25Z

	BeforeEach(func() {
		test1keyPEM, _ := pem.Decode([]byte(test1keyStr))
		if test1keyPEM == nil {
			panic(errors.New("Failed decoding private key PEM"))
		}

		test1key, err := x509.ParsePKCS1PrivateKey(test1keyPEM.Bytes)
		if err != nil {
			panic(err)
		}

		fakeapiclient = &httpclientfakes.FakeClient{}

		subject = NewDefaultManager(fakeapiclient, "fake-service", test1key)
	})

	Describe("Renew", func() {
		Context("api errors", func() {
			It("wraps api errors", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{}, errors.New("fake-err")
				}

				err := subject.Renew()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Requesting signed profile"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})

			It("wraps certificate data errors", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{
						Certificate: "invalid",
					}, nil
				}

				err := subject.Renew()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to decode PEM from certificate"))
			})

			It("wraps certificate pem errors", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{
						Certificate: test1keyStr,
					}, nil
				}

				err := subject.Renew()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing certificate"))
			})
		})

		It("works", func() {
			fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
				return api.SignUserCSRResponse{
					Certificate: test1crtValid,
					Profile:     "fake-base-profile-data",
				}, nil
			}

			err := subject.Renew()
			Expect(err).ToNot(HaveOccurred())

			profile, err := subject.GetProfile()
			Expect(err).ToNot(HaveOccurred())
			Expect(profile.BaseConfig()).To(Equal("fake-base-profile-data"))

			Expect(fakeapiclient.SignUserCSRCallCount()).To(Equal(1))

			in := fakeapiclient.SignUserCSRArgsForCall(0)

			csrPEM, _ := pem.Decode([]byte(in.CSR))
			if csrPEM == nil {
				panic("Failed to decode certificate signing request")
			}

			csr, err := x509.ParseCertificateRequest(csrPEM.Bytes)
			if err != nil {
				panic("Failed to parse certificate signing request")
			}

			Expect(csr.Subject.CommonName).To(Equal("fake-service"))
			Expect(csr.EmailAddresses).To(HaveLen(1))
			Expect(csr.EmailAddresses[0]).To(MatchRegexp(".+@.+"))
			Expect(csr.Subject.Names).To(HaveLen(2)) // 0 = CN
			Expect(csr.Subject.Names[1].Type).To(Equal(asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}))
			Expect(csr.Subject.Names[1].Value).To(Equal(csr.EmailAddresses[0]))
		})
	})

	Describe("Sign", func() {
		It("signs data", func() {
			Expect(subject.Sign([]byte("datum"))).To(Equal([]byte{121, 157, 241, 228, 204, 242, 175, 81, 96, 177, 150, 234, 233, 167, 38, 52, 15, 75, 57, 179, 40, 90, 168, 3, 215, 68, 8, 94, 71, 149, 185, 68, 225, 216, 91, 116, 203, 34, 186, 122, 230, 160, 170, 31, 175, 245, 236, 70, 228, 157, 248, 104, 66, 88, 219, 49, 253, 8, 5, 210, 166, 159, 28, 165, 111, 2, 198, 5, 221, 86, 5, 242, 175, 34, 215, 199, 135, 91, 81, 33, 247, 96, 173, 64, 203, 92, 102, 109, 191, 60, 13, 181, 89, 106, 58, 145, 118, 26, 150, 39, 160, 226, 138, 251, 41, 23, 174, 201, 17, 7, 157, 158, 50, 31, 132, 52, 134, 33, 199, 224, 82, 170, 69, 85, 205, 43, 154, 233}))
		})
	})

	Describe("GetProfile", func() {
		Context("certificate is missing", func() {
			It("renews and returns", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{
						Certificate: test1crtValid,
						Profile:     "fake-base-profile-data",
					}, nil
				}

				profile, err := subject.GetProfile()
				Expect(err).ToNot(HaveOccurred())
				Expect(profile.CertificatePEM()).To(Equal([]byte(test1crtValid)))

				Expect(fakeapiclient.SignUserCSRCallCount()).To(Equal(1))
			})

			It("wraps renewal errors errors", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{}, errors.New("fake-err")
				}

				_, err := subject.GetProfile()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Renewing certificate"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("does not try to renew valid certificates", func() {
			fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
				return api.SignUserCSRResponse{
					Certificate: test1crtValid,
					Profile:     "fake-base-profile-data",
				}, nil
			}

			err := subject.Renew()
			Expect(err).ToNot(HaveOccurred())

			Expect(fakeapiclient.SignUserCSRCallCount()).To(Equal(1))

			profile, err := subject.GetProfile()
			Expect(err).ToNot(HaveOccurred())
			Expect(profile.CertificatePEM()).To(Equal([]byte(test1crtValid)))

			Expect(fakeapiclient.SignUserCSRCallCount()).To(Equal(1))

		})
	})

	Describe("IsCertificateValid", func() {
		Context("initial manager without a certificate", func() {
			It("is not valid", func() {
				Expect(subject.IsCertificateValid()).To(BeFalse())
			})
		})

		Context("expired certificate", func() {
			It("is not valid", func() {
				fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
					return api.SignUserCSRResponse{
						Certificate: test1crtExpired,
						Profile:     "fake-base-profile-data",
					}, nil
				}

				Expect(subject.Renew()).ToNot(HaveOccurred())
				Expect(subject.IsCertificateValid()).To(BeFalse())
			})
		})

		It("is valid when not expired", func() {
			fakeapiclient.SignUserCSRStub = func(_ api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
				return api.SignUserCSRResponse{
					Certificate: test1crtValid,
					Profile:     "fake-base-profile-data",
				}, nil
			}

			Expect(subject.Renew()).ToNot(HaveOccurred())
			Expect(subject.IsCertificateValid()).To(BeTrue())
		})
	})
})
