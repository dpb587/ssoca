package req

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/certauth"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/openvpn"
	svcapi "github.com/dpb587/ssoca/service/openvpn/api"
)

type SignUserCSR struct {
	CertAuth    certauth.Provider
	Validity    time.Duration
	BaseProfile string

	req.WithAuthenticationRequired
}

var _ req.RouteHandler = SignUserCSR{}

func (h SignUserCSR) Route() string {
	return "sign-user-csr"
}

func (h SignUserCSR) Execute(request req.Request) error {
	payload := svcapi.SignUserCSRRequest{}

	err := request.ReadPayload(&payload)
	if err != nil {
		return err
	}

	response := svcapi.SignUserCSRResponse{}

	csrPEM, _ := pem.Decode([]byte(payload.CSR))
	if csrPEM == nil {
		return apierr.NewError(errors.New("Decoding CSR"), http.StatusBadRequest, "Failed to decode certificate signing request")
	}

	csr, err := x509.ParseCertificateRequest(csrPEM.Bytes)
	if err != nil {
		return apierr.NewError(errors.Wrap(err, "Parsing CSR"), http.StatusBadRequest, "Failed to parse certificate signing request")
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return errors.Wrap(err, "Generating serial number")
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{fmt.Sprintf("ssoca/%s", svc.Service{}.Version())},
			CommonName:   request.AuthToken.ID,
		},
		EmailAddresses:        csr.EmailAddresses,
		NotBefore:             now.Add(-5 * time.Second).UTC(),
		NotAfter:              now.Add(h.Validity).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certificatePEM, err := h.CertAuth.SignCertificate(&template, csr.PublicKey, request.LoggerContext)
	if err != nil {
		return errors.Wrap(err, "Signing certificate")
	}

	response.Certificate = strings.TrimSpace(string(certificatePEM))

	caCertificate, err := h.CertAuth.GetCertificatePEM()
	if err != nil {
		return errors.Wrap(err, "Loading CA certificate")
	}

	response.Profile = fmt.Sprintf(
		"%s\n<ca>\n%s\n</ca>\n",
		strings.TrimSpace(h.BaseProfile),
		strings.TrimSpace(caCertificate),
	)

	return request.WritePayload(response)
}
