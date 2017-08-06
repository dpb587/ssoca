package profile

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/dpb587/ssoca/service/openvpn/api"
	"github.com/dpb587/ssoca/service/openvpn/httpclient"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Manager struct {
	client  httpclient.Client
	service string

	privateKey *rsa.PrivateKey

	profile          string
	certificate      *x509.Certificate
	certificateBytes []byte
}

func NewManager(client httpclient.Client, service string, privateKey *rsa.PrivateKey) Manager {
	return Manager{
		client:  client,
		service: service,

		privateKey: privateKey,
	}
}

func CreateManagerAndPrivateKey(client httpclient.Client, service string) (Manager, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Manager{}, bosherr.WrapError(err, "Generating private key")
	}

	return NewManager(client, service, privateKey), nil
}

func (m Manager) Sign(data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, m.privateKey, 0, data)
}

func (m *Manager) GetProfile() (Profile, error) {
	if m.IsExpired() {
		err := m.Renew()
		if err != nil {
			return Profile{}, bosherr.WrapError(err, "Renewing certificate")
		}
	}

	return NewProfile(m.profile, m.privateKey, m.certificateBytes), nil
}

func (m Manager) IsExpired() bool {
	return m.certificate == nil || time.Now().After(m.certificate.NotAfter)
}

func (m *Manager) Renew() error {
	csrBytes, err := m.createCSR()
	if err != nil {
		return bosherr.WrapError(err, "Creating CSR")
	}

	response, err := m.client.SignUserCSR(api.SignUserCSRRequest{
		CSR: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})),
	})
	if err != nil {
		return bosherr.WrapError(err, "Requesting signed profile")
	}

	m.profile = response.Profile
	m.certificateBytes = []byte(response.Certificate)

	pem, _ := pem.Decode(m.certificateBytes)
	if pem == nil {
		return errors.New("Failed to decode PEM from certificate")
	}

	certificate, err := x509.ParseCertificate(pem.Bytes)
	if err != nil {
		panic(err)
	}

	m.certificate = certificate

	return nil
}

func (m Manager) createCSR() ([]byte, error) {
	localuser, err := user.Current()
	if err != nil {
		return nil, bosherr.WrapError(err, "Checking local user")
	}

	localhost, err := os.Hostname()
	if err != nil {
		return nil, bosherr.WrapError(err, "Checking local host")
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
		return nil, bosherr.WrapError(err, "Creating certificate request")
	}

	return csrBytes, nil
}
