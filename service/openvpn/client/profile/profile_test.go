package profile_test

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	. "github.com/dpb587/ssoca/service/openvpn/client/profile"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Profile", func() {
	var subject Profile

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
		test1keyPEM, _ := pem.Decode([]byte(test1keyStr))
		if test1keyPEM == nil {
			panic(errors.New("Failed decoding private key PEM"))
		}

		parsed, err := x509.ParsePKCS1PrivateKey(test1keyPEM.Bytes)
		if err != nil {
			panic(err)
		}

		subject = NewProfile("fake-base-profile-line1\nfake-base-profile-line2", parsed, []byte("fake-cert-data"))
	})

	Describe("CertificatePEM", func() {
		It("returns the certificate", func() {
			Expect(subject.CertificatePEM()).To(Equal([]byte("fake-cert-data")))
		})
	})

	Describe("BaseConfig", func() {
		It("returns the base config", func() {
			Expect(subject.BaseConfig()).To(Equal(`fake-base-profile-line1
fake-base-profile-line2`))
		})
	})

	Describe("StaticConfig", func() {
		It("includes base config, private key, and certificate", func() {
			Expect(subject.StaticConfig()).To(Equal(`fake-base-profile-line1
fake-base-profile-line2
<key>
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----

</key>

<cert>
fake-cert-data
</cert>
`))
		})
	})

	Describe("ManagementConfig", func() {
		It("includes base config with management directives", func() {
			Expect(subject.ManagementConfig("127.0.0.1:12345")).To(Equal(`fake-base-profile-line1
fake-base-profile-line2
management 127.0.0.1:12345
management-client

management-external-cert ssoca
management-external-key

remap-usr1 SIGHUP
`))
		})
	})
})
