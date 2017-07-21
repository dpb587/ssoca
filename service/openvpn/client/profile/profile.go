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

func (p Profile) PrivateKeyPEM() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(p.privateKey),
		},
	)
}

func (p Profile) CertificatePEM() []byte {
	return p.certificate
}

func (p Profile) BaseConfig() string {
	return p.baseConfig
}

func (p Profile) FullConfig() string {
	config := p.baseConfig

	// inline key-pair
	config = fmt.Sprintf("%s\n<key>\n%s\n</key>\n", config, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p.privateKey)}))
	config = fmt.Sprintf("%s\n<cert>\n%s\n</cert>\n", config, p.certificate)

	return config
}

func (p Profile) ManagementConfig(managementAddress string) string {
	config := p.baseConfig

	// management configuration
	config = fmt.Sprintf("%s\nmanagement %s\nmanagement-client\n", config, managementAddress)
	config = fmt.Sprintf("%s\nmanagement-external-cert ssoca\nmanagement-external-key\n", config)

	// force connection resets to flush credentials
	config = fmt.Sprintf("%s\nremap-usr1 SIGHUP\n", config)

	return config
}
