package fs

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"time"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/certauth"
)

type Provider struct {
	name   string
	config Config
	fs     boshsys.FileSystem
	logger logrus.FieldLogger
}

var _ certauth.Provider = Provider{}

func NewProvider(name string, config Config, fs boshsys.FileSystem, logger logrus.FieldLogger) Provider {
	return Provider{
		name:   name,
		config: config,
		fs:     fs,
		logger: logger,
	}
}

func (p Provider) Name() string {
	return p.name
}

func (p Provider) SignCertificate(template *x509.Certificate, publicKey interface{}, loggerContext logrus.Fields) ([]byte, error) {
	caCertificate, err := p.GetCertificate()
	if err != nil {
		return nil, errors.Wrap(err, "getting CA certificate")
	}

	caPrivateKey, err := p.getPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "getting CA private key")
	}

	certificate, err := x509.CreateCertificate(
		rand.Reader,
		template,
		caCertificate,
		publicKey,
		caPrivateKey,
	)
	if err != nil {
		return nil, errors.Wrap(err, "signing x509 certificate")
	}

	p.logger.WithFields(loggerContext).WithFields(logrus.Fields{
		"certauth.x509.serial":      template.Subject.SerialNumber,
		"certauth.x509.not_before":  template.NotBefore.Format(time.RFC3339),
		"certauth.x509.not_after":   template.NotAfter.Format(time.RFC3339),
		"certauth.x509.common_name": template.Subject.CommonName,
	}).Info("signed x509 certificate")

	bytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})

	return bytes, nil
}

func (p Provider) SignSSHCertificate(certificate *ssh.Certificate, loggerContext logrus.Fields) error {
	caPrivateKey, err := p.getPrivateKey()
	if err != nil {
		return errors.Wrap(err, "getting CA private key")
	}

	signer, err := ssh.NewSignerFromKey(caPrivateKey)
	if err != nil {
		return errors.Wrap(err, "creating ssh signer")
	}

	err = certificate.SignCert(rand.Reader, signer)
	if err != nil {
		return errors.Wrap(err, "signing ssh certificate")
	}

	p.logger.WithFields(loggerContext).WithFields(logrus.Fields{
		"certauth.ssh.valid_after":  time.Unix(int64(certificate.ValidAfter), 0).Format(time.RFC3339),
		"certauth.ssh.valid_before": time.Unix(int64(certificate.ValidBefore), 0).Format(time.RFC3339),
		"certauth.ssh.key_id":       certificate.KeyId,
	}).Info("signed ssh certificate")

	return nil
}

func (p Provider) GetCertificate() (*x509.Certificate, error) {
	str, err := p.GetCertificatePEM()
	if err != nil {
		return nil, errors.Wrap(err, "reading certificate")
	}

	certificatePEM, _ := pem.Decode([]byte(str))
	if certificatePEM == nil {
		return nil, errors.New("failed decoding certificate PEM")
	}

	certificate, err := x509.ParseCertificate(certificatePEM.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parsing certificate")
	}

	return certificate, nil
}

func (p Provider) GetCertificatePEM() (string, error) {
	return p.fs.ReadFileString(p.config.CertificatePath)
}

func (p Provider) getPrivateKey() (interface{}, error) {
	str, err := p.fs.ReadFileString(p.config.PrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading private key")
	}

	privateKeyPEM, _ := pem.Decode([]byte(str))
	if privateKeyPEM == nil {
		return nil, errors.New("failed decoding private key PEM")
	}

	privateKey, e := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if e != nil {
		return nil, errors.Wrap(e, "parsing private key")
	}

	return privateKey, nil
}
