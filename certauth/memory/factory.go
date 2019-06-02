package memory

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "loading config")
	}

	err = f.validateConfig(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "validating config")
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
		return errors.New("configuration missing: certificate")
	}

	certificatePEM, _ := pem.Decode([]byte(config.CertificateString))
	if certificatePEM == nil {
		return errors.New("failed decoding certificate PEM")
	}

	certificate, err := x509.ParseCertificate(certificatePEM.Bytes)
	if err != nil {
		return errors.Wrap(err, "parsing certificate")
	}

	config.Certificate = *certificate

	if config.PrivateKeyString == "" {
		return errors.New("configuration missing: private_key")
	}

	privateKeyPEM, _ := pem.Decode([]byte(config.PrivateKeyString))
	if privateKeyPEM == nil {
		return errors.New("failed decoding private key PEM")
	}

	privateKey, e := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if e != nil {
		return errors.Wrap(e, "parsing private key")
	}

	config.PrivateKey = privateKey

	return nil
}
