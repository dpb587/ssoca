package req

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/certauth"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BaseProfile struct {
	CertAuth    certauth.Provider
	BaseProfile string
}

func (h BaseProfile) Route() string {
	return "base-profile"
}

func (h BaseProfile) Execute(_ *auth.Token, w http.ResponseWriter) error {
	caCertificate, err := h.CertAuth.GetCertificatePEM()
	if err != nil {
		return bosherr.WrapError(err, "Loading CA certificate")
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf(
		"%s\nremap-usr1 SIGTERM\n<ca>\n%s\n</ca>\n",
		strings.TrimSpace(h.BaseProfile),
		strings.TrimSpace(caCertificate),
	)))

	return nil
}
