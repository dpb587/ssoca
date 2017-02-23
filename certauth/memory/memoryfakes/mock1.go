package memoryfakes

import (
	"github.com/dpb587/ssoca/certauth"
	. "github.com/dpb587/ssoca/certauth/memory"

	logrustest "github.com/Sirupsen/logrus/hooks/test"
)

func CreateMock1() certauth.Provider {
	logger, _ := logrustest.NewNullLogger()

	provider, err := NewFactory(logger).Create(
		"mock1",
		map[string]interface{}{
			"certificate": `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`,
			"private_key": `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		},
	)
	if err != nil {
		panic(err)
	}

	return provider
}
