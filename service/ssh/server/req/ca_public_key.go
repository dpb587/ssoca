package req

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/ssh/api"
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
		return errors.Wrap(err, "loading certificate")
	}

	sshcert, err := ssh.NewPublicKey(certificate.PublicKey)
	if err != nil {
		return errors.Wrap(err, "parsing ssh public key")
	}

	payload.OpenSSH = fmt.Sprintf("%s %s", sshcert.Type(), base64.StdEncoding.EncodeToString(sshcert.Marshal()))

	return request.WritePayload(payload)
}
