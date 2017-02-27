package req

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server"
	"github.com/dpb587/ssoca/service/ssh/api"
	"github.com/dpb587/ssoca/service/ssh/config"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type SignPublicKey struct {
	Validity        time.Duration
	CriticalOptions config.CriticalOptions
	Extensions      config.Extensions
	CertAuth        certauth.Provider
	Target          config.Target
}

func (h SignPublicKey) Route() string {
	return "sign-public-key"
}

func (h SignPublicKey) Execute(req *http.Request, token *auth.Token, payload api.SignPublicKeyRequest, loggerContext logrus.Fields) (api.SignPublicKeyResponse, error) {
	res := api.SignPublicKeyResponse{}

	parts := strings.SplitN(payload.PublicKey, " ", 3)
	if len(parts) < 2 {
		return res, server.NewAPIError(errors.New("Invalid public key format"), 400, "Failed to read public key")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return res, server.NewAPIError(bosherr.WrapErrorf(err, "Decoding public key"), 400, "Failed to decode public key")
	}

	publicKey, err := ssh.ParsePublicKey([]byte(decoded))
	if err != nil {
		return res, server.NewAPIError(bosherr.WrapErrorf(err, "Parsing public key"), 400, "Failed to parse public key")
	}

	now := time.Now()
	certificate := ssh.Certificate{
		// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.certkeys
		Key:             publicKey,
		KeyId:           token.ID,
		CertType:        ssh.UserCert,
		ValidAfter:      uint64(now.Add(-5 * time.Second).UTC().Unix()),
		ValidBefore:     uint64(now.Add(h.Validity).UTC().Unix()),
		ValidPrincipals: []string{"vcap"},
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}

	for criticalOption, criticalOptionData := range h.CriticalOptions {
		certificate.Permissions.CriticalOptions[string(criticalOption)] = criticalOptionData
	}

	for _, extension := range h.Extensions {
		if extension == config.ExtensionNoDefaults {
			continue
		}

		certificate.Permissions.Extensions[string(extension)] = ""
	}

	err = h.CertAuth.SignSSHCertificate(&certificate, loggerContext)
	if err != nil {
		return res, bosherr.WrapError(err, "Signing certificate")
	}

	res.Certificate = fmt.Sprintf("%s %s", certificate.Type(), base64.StdEncoding.EncodeToString(certificate.Marshal()))

	if h.Target.Configured() {
		res.Target = &api.SignPublicKeyTargetResponse{
			Host: h.Target.Host,
			Port: h.Target.Port,
		}

		if h.Target.User != nil {
			user, err := h.Target.User.Evaluate(req, token)
			if err != nil {
				return res, bosherr.WrapError(err, "Evaluting target.user")
			}

			res.Target.User = user
		}
	}

	return res, nil
}
