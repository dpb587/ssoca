package req

import (
	"fmt"
	"strings"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/service/req"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BaseProfile struct {
	CertAuth    certauth.Provider
	BaseProfile string

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = BaseProfile{}

func (h BaseProfile) Route() string {
	return "base-profile"
}

func (h BaseProfile) Execute(request req.Request) error {
	caCertificate, err := h.CertAuth.GetCertificatePEM()
	if err != nil {
		return bosherr.WrapError(err, "Loading CA certificate")
	}

	request.RawResponse.Header().Add("Content-Type", "text/plain")
	request.RawResponse.Write([]byte(fmt.Sprintf(
		"%s\nremap-usr1 SIGTERM\n<ca>\n%s\n</ca>\n",
		strings.TrimSpace(h.BaseProfile),
		strings.TrimSpace(caCertificate),
	)))

	return nil
}
