package profile

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Profile struct {
	baseConfig  string
	privateKey  *rsa.PrivateKey
	certificate []byte
}

func NewProfile(baseConfig string, privateKey *rsa.PrivateKey, certificate []byte) Profile {
	return Profile{
		baseConfig:  baseConfig,
		privateKey:  privateKey,
		certificate: certificate,
	}
}

func (p Profile) CertificatePEM() []byte {
	return p.certificate
}

func (p Profile) BaseConfig() string {
	return p.baseConfig
}

func (p Profile) StaticConfig() string {
	config := p.BaseConfig()

	// with static configs, short-lived certificates may expire or be rejected and
	// openvpn may continue retrying; on the assumption that the process manager
	// will restart and ssoca will generate a new, valid certificate, simply exit.
	config = fmt.Sprintf("%s\nremap-usr1 SIGTERM\n", config)

	// inline key-pair
	config = fmt.Sprintf("%s\n<key>\n%s\n</key>\n", config, p.privateKeyPEM())
	config = fmt.Sprintf("%s\n<cert>\n%s\n</cert>\n", config, p.certificate)

	return config
}

func (p Profile) ManagementConfig(managementAddress, managementPasswordFile string) string {
	config := p.BaseConfig()

	config = config + "\n"
	config = fmt.Sprintf("%smanagement %s %s\n", config, managementAddress, managementPasswordFile)
	config = fmt.Sprintf("%smanagement-client\n", config)
	config = fmt.Sprintf("%smanagement-external-cert ssoca\n", config)
	config = fmt.Sprintf("%smanagement-external-key\n", config)

	// since openvpn does not know how to receive certificate updates, when USR1
	// connection-reset is encountered, have it reload the configuration and
	// management connection to receive newer certificates.
	config = fmt.Sprintf("%s\nremap-usr1 SIGHUP\n", config)

	return config
}

func (p Profile) privateKeyPEM() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(p.privateKey),
		},
	)
}
