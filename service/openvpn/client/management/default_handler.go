package management

import (
	"encoding/base64"
	"io"

	"github.com/dpb587/ssoca/service/openvpn/client/profile"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DefaultHandler struct {
	profileManager profile.Manager
}

var _ ServerHandler = &DefaultHandler{}

func NewDefaultHandler(profileManager profile.Manager) *DefaultHandler {
	return &DefaultHandler{
		profileManager: profileManager,
	}
}

func (ch *DefaultHandler) NeedCertificate(w io.Writer, _ string) (ServerHandlerCallback, error) {
	profile, err := ch.profileManager.GetProfile()
	if err != nil {
		return nil, bosherr.WrapError(err, "Retrieving profile")
	}

	w.Write([]byte("certificate\n"))
	w.Write(profile.CertificatePEM())
	w.Write([]byte("\n"))
	w.Write([]byte("END\n"))

	return SuccessCallback, nil
}

func (ch *DefaultHandler) SignRSA(w io.Writer, data string) (ServerHandlerCallback, error) {
	if !ch.profileManager.IsCertificateValid() {
		w.Write([]byte("signal SIGHUP\n"))

		return SuccessCallback, nil
	}

	data64, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, bosherr.WrapError(err, "Decoding signing token")
	}

	signature, err := ch.profileManager.Sign(data64)
	if err != nil {
		return nil, bosherr.WrapError(err, "Signing token")
	}

	signature64 := base64.StdEncoding.EncodeToString(signature)

	w.Write([]byte("rsa-sig\n"))
	w.Write([]byte(signature64))
	w.Write([]byte("\n"))
	w.Write([]byte("END\n"))

	return SuccessCallback, nil
}
