package config

import (
	"crypto"
	"crypto/x509"
	"time"
)

type Config struct {
	IDPMetadata string        `yaml:"idp_metadata"`
	KeyPair     KeyPairConfig `yaml:"key_pair"`

	JWT JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	PrivateKey   string        `yaml:"private_key"`
	Validity     time.Duration `yaml:"validity"`
	ValidityPast time.Duration `yaml:"validity_past"`
}

type KeyPairConfig struct {
	CertificateString string `yaml:"certificate"`
	PrivateKeyString  string `yaml:"private_key"`

	Certificate x509.Certificate  `yaml:"-"`
	PrivateKey  crypto.PrivateKey `yaml:"-"`
}
