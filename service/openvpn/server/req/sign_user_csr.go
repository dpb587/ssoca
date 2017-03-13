package req

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/api"
	"github.com/dpb587/ssoca/server/service/req"
	svc "github.com/dpb587/ssoca/service/openvpn"
	svcapi "github.com/dpb587/ssoca/service/openvpn/api"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type SignUserCSR struct {
	CertAuth    certauth.Provider
	Validity    time.Duration
	BaseProfile string
}

var _ req.RouteHandler = SignUserCSR{}

func (h SignUserCSR) Route() string {
	return "sign-user-csr"
}

func (h SignUserCSR) Execute(token *auth.Token, payload svcapi.SignUserCSRRequest, loggerContext logrus.Fields) (svcapi.SignUserCSRResponse, error) {
	res := svcapi.SignUserCSRResponse{}

	csrPEM, _ := pem.Decode([]byte(payload.CSR))
	if csrPEM == nil {
		return res, api.NewError(errors.New("Decoding CSR"), http.StatusBadRequest, "Failed to decode certificate signing request")
	}

	csr, err := x509.ParseCertificateRequest(csrPEM.Bytes)
	if err != nil {
		return res, api.NewError(bosherr.WrapError(err, "Parsing CSR"), http.StatusBadRequest, "Failed to parse certificate signing request")
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return res, bosherr.WrapError(err, "Generating serial number")
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{fmt.Sprintf("ssoca/%s", svc.Service{}.Version())},
			CommonName:   token.ID,
		},
		EmailAddresses:        csr.EmailAddresses,
		NotBefore:             now.Add(-5 * time.Second).UTC(),
		NotAfter:              now.Add(h.Validity).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certificatePEM, err := h.CertAuth.SignCertificate(&template, csr.PublicKey, loggerContext)
	if err != nil {
		return res, bosherr.WrapError(err, "Signing certificate")
	}

	res.Certificate = strings.TrimSpace(string(certificatePEM))

	caCertificate, err := h.CertAuth.GetCertificatePEM()
	if err != nil {
		return res, bosherr.WrapError(err, "Loading CA certificate")
	}

	res.Profile = fmt.Sprintf(
		"%s\nremap-usr1 SIGTERM\n<ca>\n%s\n</ca>\n<cert>\n%s\n</cert>\n",
		strings.TrimSpace(h.BaseProfile),
		strings.TrimSpace(caCertificate),
		res.Certificate,
	)

	return res, nil
}
