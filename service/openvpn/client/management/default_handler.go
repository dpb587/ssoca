package management

import (
	"encoding/base64"
	"io"

	"github.com/dpb587/ssoca/service/openvpn/client/profile"
)

type DefaultHandler struct {
	profileManager *profile.Manager
}

var _ ClientHandler = &DefaultHandler{}

func NewDefaultHandler(profileManager *profile.Manager) *DefaultHandler {
	return &DefaultHandler{
		profileManager: profileManager,
	}
}

func (ch *DefaultHandler) NeedCertificate(w io.Writer, _ string) (ClientHandlerCallback, error) {
	profile, err := ch.profileManager.GetProfile()
	if err != nil {
		return nil, err
	}

	w.Write([]byte("certificate\n"))
	w.Write(profile.CertificatePEM())
	w.Write([]byte("\n"))
	w.Write([]byte("END\n"))

	return SimpleCallbackHandler, nil
}

func (ch *DefaultHandler) SignRSA(w io.Writer, data string) (ClientHandlerCallback, error) {
	if ch.profileManager.IsExpired() {
		w.Write([]byte("signal SIGHUP\n"))

		return SimpleCallbackHandler, nil
	}

	data64, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	signature, err := ch.profileManager.Sign(data64)
	if err != nil {
		return nil, err
	}

	signature64 := base64.StdEncoding.EncodeToString(signature)

	w.Write([]byte("rsa-sig\n"))
	w.Write([]byte(signature64))
	w.Write([]byte("\n"))
	w.Write([]byte("END\n"))

	return SimpleCallbackHandler, nil
}
