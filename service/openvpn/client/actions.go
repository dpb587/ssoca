package client

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
	"strings"

	"github.com/dpb587/ssoca/service/openvpn/api"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (s Service) CreateEphemeralKey(service string) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Generating private key")
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	localuser, err := user.Current()
	if err != nil {
		return "", "", bosherr.WrapError(err, "Checking local user")
	}

	localhost, err := os.Hostname()
	if err != nil {
		return "", "", bosherr.WrapError(err, "Checking local host")
	}

	emailAddress := fmt.Sprintf("%s@%s", localuser.Username, localhost)

	subj := pkix.Name{CommonName: service}

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

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Creating certificate request")
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	return strings.TrimSpace(string(privateKeyPEM)), strings.TrimSpace(string(csrPEM)), nil
}

func (s Service) CreateProfile(service string) (string, error) {
	privateKey, csr, err := s.CreateEphemeralKey(service)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating key pair")
	}

	client, err := s.GetClient(service)
	if err != nil {
		return "", bosherr.WrapError(err, "Getting client")
	}

	response, err := client.SignUserCSR(api.SignUserCSRRequest{
		CSR: csr,
	})
	if err != nil {
		return "", bosherr.WrapError(err, "Requesting signed profile")
	}

	return fmt.Sprintf("%s<key>\n%s\n</key>\n", response.Profile, privateKey), nil
}
