package memory_test

import (
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"

	. "github.com/dpb587/ssoca/certauth/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	var subject Factory
	var logger logrus.FieldLogger

	// certstrap init --key-bits 1024 --common-name ssoca-test --passphrase ''
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
	var ca1keyStr = `-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDpN3e+wD9/2UdA94jMkH0nUlNdoNV+qgKVTZlGW5rsGaV9c5Ce
Cj/U6Z607dMcSv6bEUwaB8lnoK6MFF3cDv/eH4jDvaMoYFqDiIQEj24FzJ5GBZ3N
xXuXt3NBPTRcIGeQklElXiOgMii+q4DPqIp/gCXjDJDmTVEz0oCVSKkgWQIDAQAB
AoGBANC3T3drXmjw74/4+Hj7Jsa2Kt20Pt1pEX7FP9Nz0CZUnYK0lkyaJ55IpjyO
S00a4NmulUkGhv0zFINRBt8WnW1bjBxNmqyBYh2diO3vA/gk8U1gcifW1LQt8WmE
ietvN3OFXI1a7FipchCZYQn5Rr8O3a/tjwohtWIDdaDltw+xAkEA7Ybxu8OXQnvy
Y+fDISRGG5vDFGnNGe9KcREIxSF6LWJ7+ap5LmMxnhfag5qlrObQW3K2miTpGYkl
CIRRNFMIvwJBAPtatE1evu25R3NSTU2YwQgkEymh40PW+lncYge6ZqZGfK7J5JBK
wr1ug7KjTJgIfY2Sg2VHn56HAdA4RUl2xOcCQQDZqnTxpQ6DHYSFqwg04cHhYP8H
QOF0Z8WnEX4g8Em/N2X26BK+wKXig2d6fIhghu/fLaNKZJK8FOK8CE1GDuWPAkEA
wrP6Ysx3vZH+JPil5Ovk6zd2mJNMhmpqt10dmrrrdPW483R01sjynOaUobYZSNOa
3iWWHsgifxw5bV+JXGTiFQJBAKwh6Hvli5hcfoepPMz2RQnmU1NM8hJOHHeZh+eT
z6hlMpOS9rSjABcBdXxXjFXtIEjWUG5Tj8yOYd735zY8Ny8=
-----END RSA PRIVATE KEY-----`

	BeforeEach(func() {
		logger, _ = logrustest.NewNullLogger()

		subject = NewFactory(logger)
	})

	Describe("Create", func() {
		It("works", func() {
			provider, err := subject.Create("name1", map[string]interface{}{
				"certificate": ca1crtStr,
				"private_key": ca1keyStr,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(provider.Name()).To(Equal("name1"))
		})

		Context("invalid yaml", func() {
			It("remarshals configuration", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": map[string]interface{}{},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Loading config"))
			})
		})

		Context("certificate/key errors", func() {
			It("errors on missing certificate", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"private_key": ca1keyStr,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Configuration missing: certificate"))
			})

			It("errors on misconfigured certificate", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": ca1keyStr,
					"private_key": ca1keyStr,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing certificate"))
			})

			It("errors on invalid certificate", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": "broken",
					"private_key": ca1keyStr,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed decoding certificate"))
			})

			It("errors on missing private key", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": ca1crtStr,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Configuration missing: private_key"))
			})

			It("errors on misconfigured private key", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": ca1crtStr,
					"private_key": ca1crtStr,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing private key"))
			})

			It("errors on invalid private key", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate": ca1crtStr,
					"private_key": "broken",
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed decoding private key"))
			})
		})
	})
})
