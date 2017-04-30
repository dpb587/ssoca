package profile

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Settings struct {
	Profile     string
	PrivateKey  *rsa.PrivateKey
	Certificate []byte
}

func (s Settings) PrivateKeyPEM() []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(s.PrivateKey),
		},
	)
}

func (s Settings) String() string {
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(s.PrivateKey),
		},
	)

	return fmt.Sprintf("%s\n<key>\n%s\n</key>\n<cert>\n%s\n</cert>\n", s.Profile, privateKeyPEM, s.Certificate)
}
