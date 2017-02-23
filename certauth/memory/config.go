package memory

import (
	"crypto"
	"crypto/x509"
)

type Config struct {
	CertificateString string `yaml:"certificate"`
	PrivateKeyString  string `yaml:"private_key"`

	Certificate x509.Certificate  `yaml:"-"`
	PrivateKey  crypto.PrivateKey `yaml:"-"`
}
