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
	"github.com/dpb587/ssoca/server/api"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
	svcapi "github.com/dpb587/ssoca/service/ssh/api"
	svcconfig "github.com/dpb587/ssoca/service/ssh/config"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type SignPublicKey struct {
	Validity        time.Duration
	Principals      []dynamicvalue.Value
	CriticalOptions svcconfig.CriticalOptions
	Extensions      svcconfig.Extensions
	CertAuth        certauth.Provider
	Target          svcconfig.Target
}

func (h SignPublicKey) Route() string {
	return "sign-public-key"
}

func (h SignPublicKey) Execute(req *http.Request, token *auth.Token, payload svcapi.SignPublicKeyRequest, loggerContext logrus.Fields) (svcapi.SignPublicKeyResponse, error) {
	res := svcapi.SignPublicKeyResponse{}

	parts := strings.SplitN(payload.PublicKey, " ", 3)
	if len(parts) < 2 {
		return res, api.NewError(errors.New("Invalid public key format"), 400, "Failed to read public key")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return res, api.NewError(bosherr.WrapErrorf(err, "Decoding public key"), 400, "Failed to decode public key")
	}

	publicKey, err := ssh.ParsePublicKey([]byte(decoded))
	if err != nil {
		return res, api.NewError(bosherr.WrapErrorf(err, "Parsing public key"), 400, "Failed to parse public key")
	}

	now := time.Now()
	certificate := ssh.Certificate{
		// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.certkeys
		Key:             publicKey,
		KeyId:           token.ID,
		CertType:        ssh.UserCert,
		ValidAfter:      uint64(now.Add(-5 * time.Second).UTC().Unix()),
		ValidBefore:     uint64(now.Add(h.Validity).UTC().Unix()),
		ValidPrincipals: []string{},
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}

	for _, dynamicPrincipal := range h.Principals {
		principal, err1 := dynamicPrincipal.Evaluate(req, token)
		if err1 != nil {
			return res, bosherr.WrapError(err1, "Evaluating dynamic principal")
		} else if principal == "" {
			continue
		}

		certificate.ValidPrincipals = append(certificate.ValidPrincipals, principal)
	}

	for criticalOption, criticalOptionData := range h.CriticalOptions {
		certificate.Permissions.CriticalOptions[string(criticalOption)] = criticalOptionData
	}

	for _, extension := range h.Extensions {
		certificate.Permissions.Extensions[string(extension)] = ""
	}

	err = h.CertAuth.SignSSHCertificate(&certificate, loggerContext)
	if err != nil {
		return res, bosherr.WrapError(err, "Signing certificate")
	}

	res.Certificate = fmt.Sprintf("%s %s", certificate.Type(), base64.StdEncoding.EncodeToString(certificate.Marshal()))

	if h.Target.Configured() {
		res.Target = &svcapi.SignPublicKeyTargetResponse{
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
