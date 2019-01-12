package profile

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/openvpn/api"
	"github.com/dpb587/ssoca/service/openvpn/httpclient"
)

type DefaultManager struct {
	client  httpclient.Client
	service string

	privateKey *rsa.PrivateKey

	profile          string
	certificate      *x509.Certificate
	certificateBytes []byte
}

var _ Manager = &DefaultManager{}

func NewDefaultManager(client httpclient.Client, service string, privateKey *rsa.PrivateKey) DefaultManager {
	return DefaultManager{
		client:  client,
		service: service,

		privateKey: privateKey,
	}
}

func CreateManagerAndPrivateKey(client httpclient.Client, service string) (DefaultManager, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return DefaultManager{}, errors.Wrap(err, "Generating private key")
	}

	return NewDefaultManager(client, service, privateKey), nil
}

func (m DefaultManager) Sign(data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, m.privateKey, 0, data)
}

func (m *DefaultManager) GetProfile() (Profile, error) {
	if !m.IsCertificateValid() {
		err := m.Renew()
		if err != nil {
			return Profile{}, errors.Wrap(err, "Renewing certificate")
		}
	}

	return NewProfile(m.profile, m.privateKey, m.certificateBytes), nil
}

func (m DefaultManager) IsCertificateValid() bool {
	return m.certificate != nil && time.Now().Before(m.certificate.NotAfter)
}

func (m *DefaultManager) Renew() error {
	csrBytes, err := m.createCSR()
	if err != nil {
		return errors.Wrap(err, "Creating CSR")
	}

	response, err := m.client.SignUserCSR(api.SignUserCSRRequest{
		CSR: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})),
	})
	if err != nil {
		return errors.Wrap(err, "Requesting signed profile")
	}

	m.profile = response.Profile
	m.certificateBytes = []byte(response.Certificate)

	pem, _ := pem.Decode(m.certificateBytes)
	if pem == nil {
		return errors.New("Failed to decode PEM from certificate")
	}

	certificate, err := x509.ParseCertificate(pem.Bytes)
	if err != nil {
		return errors.Wrap(err, "Parsing certificate")
	}

	m.certificate = certificate

	return nil
}

func (m DefaultManager) createCSR() ([]byte, error) {
	localuser, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Checking local user")
	}

	localhost, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "Checking local hostname")
	}

	emailAddress := fmt.Sprintf("%s@%s", localuser.Username, localhost)

	subj := pkix.Name{CommonName: m.service}

	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{
			Type:  asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1},
			Value: emailAddress,
		},
	})

	asn1Subj, _ := asn1.Marshal(rawSubj)
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, m.privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Creating certificate request")
	}

	return csrBytes, nil
}
