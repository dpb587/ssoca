package memory

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/certauth"
)

type Provider struct {
	name   string
	config Config
	logger logrus.FieldLogger
}

var _ certauth.Provider = Provider{}

func NewProvider(name string, config Config, logger logrus.FieldLogger) Provider {
	return Provider{
		name:   name,
		config: config,
		logger: logger,
	}
}

func (p Provider) Name() string {
	return p.name
}

func (p Provider) SignCertificate(template *x509.Certificate, publicKey interface{}, loggerContext logrus.Fields) ([]byte, error) {
	caCertificate := p.config.Certificate
	caPrivateKey := p.config.PrivateKey

	certificate, err := x509.CreateCertificate(
		rand.Reader,
		template,
		&caCertificate,
		publicKey,
		caPrivateKey,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Signing x509 certificate")
	}

	p.logger.WithFields(loggerContext).WithFields(logrus.Fields{
		"certauth.x509.serial":      template.Subject.SerialNumber,
		"certauth.x509.not_before":  template.NotBefore.Format(time.RFC3339),
		"certauth.x509.not_after":   template.NotAfter.Format(time.RFC3339),
		"certauth.x509.common_name": template.Subject.CommonName,
	}).Info("Signed x509 certificate")

	bytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})

	return bytes, nil
}

func (p Provider) SignSSHCertificate(certificate *ssh.Certificate, loggerContext logrus.Fields) error {
	caPrivateKey := p.config.PrivateKey

	signer, err := ssh.NewSignerFromKey(caPrivateKey)
	if err != nil {
		return errors.Wrap(err, "Creating ssh signer")
	}

	err = certificate.SignCert(rand.Reader, signer)
	if err != nil {
		return errors.Wrap(err, "Signing ssh certificate")
	}

	p.logger.WithFields(loggerContext).WithFields(logrus.Fields{
		"certauth.ssh.valid_after":  time.Unix(int64(certificate.ValidAfter), 0).Format(time.RFC3339),
		"certauth.ssh.valid_before": time.Unix(int64(certificate.ValidBefore), 0).Format(time.RFC3339),
		"certauth.ssh.key_id":       certificate.KeyId,
	}).Info("Signed ssh certificate")

	return nil
}

func (p Provider) GetCertificate() (*x509.Certificate, error) {
	return &p.config.Certificate, nil
}

func (p Provider) GetCertificatePEM() (string, error) {
	return p.config.CertificateString, nil
}
