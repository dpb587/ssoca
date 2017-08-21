package memory

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/sirupsen/logrus"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config"
)

type Factory struct {
	logger logrus.FieldLogger
}

var _ certauth.ProviderFactory = Factory{}

func NewFactory(logger logrus.FieldLogger) Factory {
	return Factory{
		logger: logger,
	}
}

func (f Factory) Create(name string, options map[string]interface{}) (certauth.Provider, error) {
	var cfg Config

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	err = f.validateConfig(&cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Validating config")
	}

	provider := NewProvider(
		name,
		cfg,
		f.logger.WithFields(logrus.Fields{
			"certauth.name": name,
		}),
	)

	return provider, nil
}

func (f Factory) validateConfig(config *Config) error {
	if config.CertificateString == "" {
		return errors.New("Configuration missing: certificate")
	}

	certificatePEM, _ := pem.Decode([]byte(config.CertificateString))
	if certificatePEM == nil {
		return errors.New("Failed decoding certificate PEM")
	}

	certificate, err := x509.ParseCertificate(certificatePEM.Bytes)
	if err != nil {
		return bosherr.WrapError(err, "Parsing certificate")
	}

	config.Certificate = *certificate

	if config.PrivateKeyString == "" {
		return errors.New("Configuration missing: private_key")
	}

	privateKeyPEM, _ := pem.Decode([]byte(config.PrivateKeyString))
	if privateKeyPEM == nil {
		return errors.New("Failed decoding private key PEM")
	}

	privateKey, e := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if e != nil {
		return bosherr.WrapError(e, "Parsing private key")
	}

	config.PrivateKey = privateKey

	return nil
}
