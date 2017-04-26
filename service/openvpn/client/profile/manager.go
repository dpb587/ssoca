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

	"github.com/dpb587/ssoca/service/openvpn/api"
	"github.com/dpb587/ssoca/service/openvpn/httpclient"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Manager struct {
	client     *httpclient.Client
	service    string
	privateKey *rsa.PrivateKey
}

func NewManager(client *httpclient.Client, service string, privateKey *rsa.PrivateKey) Manager {
	return Manager{
		client:     client,
		service:    service,
		privateKey: privateKey,
	}
}

func CreateManagerAndPrivateKey(client *httpclient.Client, service string) (Manager, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return Manager{}, bosherr.WrapError(err, "Generating private key")
	}

	return NewManager(client, service, privateKey), nil
}

func (m Manager) Sign(data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, m.privateKey, 0, data)
}

func (m Manager) GetProfile() (Settings, error) {
	csrBytes, err := m.createCSR()
	if err != nil {
		return Settings{}, bosherr.WrapError(err, "Creating csr")
	}

	response, err := m.client.SignUserCSR(api.SignUserCSRRequest{
		CSR: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})),
	})
	if err != nil {
		return Settings{}, bosherr.WrapError(err, "Requesting signed profile")
	}

	return Settings{
		Profile:     response.Profile,
		PrivateKey:  m.privateKey,
		Certificate: []byte(response.Certificate),
	}, nil
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
