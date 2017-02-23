package req_test

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/dpb587/ssoca/certauth/certauthfakes"
	. "github.com/dpb587/ssoca/service/ssh/server/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaPublicKey", func() {
	var subject CAPublicKey

	Describe("Route", func() {
		It("returns", func() {
			Expect(subject.Route()).To(Equal("ca-public-key"))
		})
	})

	Describe("Execute", func() {
		var certauth certauthfakes.FakeProvider
		var ca1crtStr = `-----BEGIN CERTIFICATE-----
MIIB5TCCAU6gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwpzc29j
YS10ZXN0MB4XDTE3MDIxMzIwMzMwOFoXDTI3MDIxMzIwMzMwOFowFTETMBEGA1UE
AxMKc3NvY2EtdGVzdDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA6Td3vsA/
f9lHQPeIzJB9J1JTXaDVfqoClU2ZRlua7BmlfXOQngo/1OmetO3THEr+mxFMGgfJ
Z6CujBRd3A7/3h+Iw72jKGBag4iEBI9uBcyeRgWdzcV7l7dzQT00XCBnkJJRJV4j
oDIovquAz6iKf4Al4wyQ5k1RM9KAlUipIFkCAwEAAaNFMEMwDgYDVR0PAQH/BAQD
AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFP8lIbNl3zZPEHF17cFU
NFsK/0/oMA0GCSqGSIb3DQEBCwUAA4GBADMCd4nzc19voa60lNknhsihcfyNUeUt
EEsLCceK+9F1u2Xdj+mTNOh3MI+5m7wmFLiHuUtQovHMJ4xUpoHa6Iznc+QCbow4
SMO3sf1847tASv3eUFwEUt9vv39vtey6C6ftiUUImzZYfx6FO/A62uGEg2w3IOJ+
3cCXYiulfsyv
-----END CERTIFICATE-----`

		pemToCertificate := func(bytes []byte) *x509.Certificate {
			pem, _ := pem.Decode(bytes)
			if pem == nil {
				panic("Failed decoding PEM")
			}

			certificate, err := x509.ParseCertificate(pem.Bytes)
			if err != nil {
				panic(err)
			}

			return certificate
		}

		BeforeEach(func() {
			certauth = certauthfakes.FakeProvider{}

			subject = CAPublicKey{
				CertAuth: &certauth,
			}
		})

		It("works", func() {
			certauth.GetCertificateReturns(pemToCertificate([]byte(ca1crtStr)), nil)

			res, err := subject.Execute()

			Expect(err).ToNot(HaveOccurred())
			Expect(res.OpenSSH).To(Equal("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDpN3e+wD9/2UdA94jMkH0nUlNdoNV+qgKVTZlGW5rsGaV9c5CeCj/U6Z607dMcSv6bEUwaB8lnoK6MFF3cDv/eH4jDvaMoYFqDiIQEj24FzJ5GBZ3NxXuXt3NBPTRcIGeQklElXiOgMii+q4DPqIp/gCXjDJDmTVEz0oCVSKkgWQ=="))
		})

		Context("certauth errors", func() {
			It("errors", func() {
				certauth.GetCertificateReturns(nil, errors.New("fake-err"))

				_, err := subject.Execute()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err"))
				Expect(err.Error()).To(ContainSubstring("Loading certificate"))
			})
		})
	})
})
