package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"

	svc "github.com/dpb587/ssoca/auth/authn/saml"
	svcconfig "github.com/dpb587/ssoca/auth/authn/saml/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/server/service"
)

type ServiceFactory struct {
	endpointURL string
	failureURL  string
	successURL  string
}

func NewServiceFactory(endpointURL string, failureURL string, successURL string) ServiceFactory {
	return ServiceFactory{
		endpointURL: endpointURL,
		failureURL:  failureURL,
		successURL:  successURL,
	}
}
func (f ServiceFactory) Type() string {
	return svc.Service{}.Type()
}

func (f ServiceFactory) Create(name string, options map[string]interface{}) (service.Service, error) {
	var cfg svcconfig.Config
	cfg.JWT.Validity = 24 * time.Hour
	cfg.JWT.ValidityPast = 2 * time.Second

	err := config.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading config")
	}

	certificatePEM, _ := pem.Decode([]byte(cfg.KeyPair.CertificateString))
	if certificatePEM == nil {
		return nil, errors.New("Failed decoding certificate PEM")
	}

	certificate, err := x509.ParseCertificate(certificatePEM.Bytes)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing certificate")
	}

	cfg.KeyPair.Certificate = *certificate

	if cfg.KeyPair.PrivateKeyString == "" {
		return nil, errors.New("Configuration missing: private_key")
	}

	privateKeyPEM, _ := pem.Decode([]byte(cfg.KeyPair.PrivateKeyString))
	if privateKeyPEM == nil {
		return nil, errors.New("Failed decoding private key PEM")
	}

	privateKey, e := x509.ParsePKCS1PrivateKey(privateKeyPEM.Bytes)
	if e != nil {
		return nil, bosherr.WrapError(e, "Parsing private key")
	}

	cfg.KeyPair.PrivateKey = privateKey

	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.JWT.PrivateKey))
	// if err != nil {
	// 	return nil, bosherr.WrapError(err, "Parsing private key")
	// }

	idpEntityDescriptor := &saml.Metadata{}

	err = xml.Unmarshal([]byte(cfg.IDPMetadata), idpEntityDescriptor)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing IDP entity descriptor")
	}

	acsURL, err := url.Parse(fmt.Sprintf("%s/%s/acs", f.endpointURL, name))
	if err != nil {
		return nil, bosherr.WrapError(err, "Generating ACS endpoint")
	}

	metadataURL, err := url.Parse(fmt.Sprintf("%s/%s/metadata", f.endpointURL, name))
	if err != nil {
		return nil, bosherr.WrapError(err, "Generating ACS endpoint")
	}

	sp := &saml.ServiceProvider{
		Key:         cfg.KeyPair.PrivateKey.(*rsa.PrivateKey),
		Logger:      logger.DefaultLogger,
		Certificate: &cfg.KeyPair.Certificate,
		AcsURL:      *acsURL,
		MetadataURL: *metadataURL,
		IDPMetadata: idpEntityDescriptor,
	}

	return NewService(name, cfg, sp), nil
}
