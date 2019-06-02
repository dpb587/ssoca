package req

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/certauth"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
	"github.com/dpb587/ssoca/server/service/req"
	svcapi "github.com/dpb587/ssoca/service/ssh/api"
	svcconfig "github.com/dpb587/ssoca/service/ssh/server/config"
)

type SignPublicKey struct {
	Validity        time.Duration
	Principals      dynamicvalue.MultiValue
	CriticalOptions svcconfig.CriticalOptions
	Extensions      svcconfig.Extensions
	CertAuth        certauth.Provider
	Target          svcconfig.Target

	req.WithAuthenticationRequired
}

var _ req.RouteHandler = SignPublicKey{}

func (h SignPublicKey) Route() string {
	return "sign-public-key"
}

func (h SignPublicKey) Execute(request req.Request) error {
	response := svcapi.SignPublicKeyResponse{}

	payload := svcapi.SignPublicKeyRequest{}

	err := request.ReadPayload(&payload)
	if err != nil {
		return err
	}

	parts := strings.SplitN(payload.PublicKey, " ", 3)
	if len(parts) < 2 {
		return apierr.NewError(errors.New("invalid public key format"), 400, "failed to read public key")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return apierr.NewError(errors.Wrapf(err, "decoding public key"), 400, "failed to decode public key")
	}

	publicKey, err := ssh.ParsePublicKey([]byte(decoded))
	if err != nil {
		return apierr.NewError(errors.Wrapf(err, "parsing public key"), 400, "failed to parse public key")
	}

	now := time.Now()
	certificate := ssh.Certificate{
		// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.certkeys
		Key:             publicKey,
		KeyId:           request.AuthToken.ID,
		CertType:        ssh.UserCert,
		ValidAfter:      uint64(now.Add(-5 * time.Second).UTC().Unix()),
		ValidBefore:     uint64(now.Add(h.Validity).UTC().Unix()),
		ValidPrincipals: []string{},
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}

	principals, err := h.Principals.Evaluate(request.RawRequest, request.AuthToken)
	if err != nil {
		return errors.Wrap(err, "evaluating principals")
	}

	principalsFiltered := []string{}

	for _, principalsCandidate := range principals {
		if principalsCandidate == "" {
			continue
		}

		principalsFiltered = append(principalsFiltered, principalsCandidate)
	}

	certificate.ValidPrincipals = principalsFiltered

	criticalOptions, err := h.CriticalOptions.Evaluate(request.RawRequest, request.AuthToken)
	if err != nil {
		return errors.Wrap(err, "evaluating critical options")
	}

	for criticalOption, criticalOptionData := range criticalOptions {
		if criticalOptionData == "" {
			continue
		}

		certificate.Permissions.CriticalOptions[string(criticalOption)] = criticalOptionData
	}

	for _, extension := range h.Extensions {
		certificate.Permissions.Extensions[string(extension)] = ""
	}

	err = h.CertAuth.SignSSHCertificate(&certificate, request.LoggerContext)
	if err != nil {
		return errors.Wrap(err, "signing certificate")
	}

	response.Certificate = fmt.Sprintf("%s %s", certificate.Type(), base64.StdEncoding.EncodeToString(certificate.Marshal()))

	{
		target := &svcapi.SignPublicKeyTargetResponse{
			Host:      h.Target.Host,
			Port:      h.Target.Port,
			PublicKey: h.Target.PublicKey,
		}

		targetUser, err := h.Target.User.Evaluate(request.RawRequest, request.AuthToken)
		if err != nil {
			return errors.Wrap(err, "evaluting target user")
		}

		target.User = targetUser

		if target.Host != "" || target.Port != 0 || target.User != "" {
			response.Target = target
		}
	}

	return request.WritePayload(response)
}
