package req

import (
	"encoding/base64"
	"fmt"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/ssh/api"
	"golang.org/x/crypto/ssh"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CAPublicKey struct {
	CertAuth certauth.Provider

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = CAPublicKey{}

func (CAPublicKey) Route() string {
	return "ca-public-key"
}

func (h CAPublicKey) Execute(request req.Request) error {
	payload := api.CAPublicKeyResponse{}

	certificate, err := h.CertAuth.GetCertificate()
	if err != nil {
		return bosherr.WrapError(err, "Loading certificate")
	}

	sshcert, err := ssh.NewPublicKey(certificate.PublicKey)
	if err != nil {
		return bosherr.WrapError(err, "Parsing ssh public key")
	}

	payload.OpenSSH = fmt.Sprintf("%s %s", sshcert.Type(), base64.StdEncoding.EncodeToString(sshcert.Marshal()))

	return request.WritePayload(payload)
}
