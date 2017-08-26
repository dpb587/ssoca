package management

import (
	"encoding/base64"
	"io"
	"time"

	"github.com/dpb587/ssoca/service/openvpn/client/profile"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DefaultHandler struct {
	profileManager profile.Manager

	attempts       []time.Time
	attemptsOffset int
	attemptsMax    int
}

var _ ServerHandler = &DefaultHandler{}

func NewDefaultHandler(profileManager profile.Manager) *DefaultHandler {
	attemptsMax := 5

	return &DefaultHandler{
		profileManager: profileManager,
		attempts:       make([]time.Time, attemptsMax),
		attemptsMax:    attemptsMax,
	}
}

func (ch *DefaultHandler) NeedCertificate(w io.Writer, _ string) (ServerHandlerCallback, error) {
	profile, err := ch.profileManager.GetProfile()
	if err != nil {
		return nil, bosherr.WrapError(err, "Retrieving profile")
	} else if ch.exceedsRecentAttempts() {
		// internally throttle ourselves after repeated certificate lookups
		// typically this would mean non-expired certs are expired/invalid and should be reloaded
		w.Write([]byte("signal SIGHUP\n"))

		return SuccessCallback, nil
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

func (ch *DefaultHandler) exceedsRecentAttempts() bool {
	ch.attemptsOffset = (ch.attemptsOffset + 1) % ch.attemptsMax
	ch.attempts[ch.attemptsOffset] = time.Now()

	since := time.Now().Add(-5 * time.Minute)
	sinceCount := 0

	for _, attempt := range ch.attempts {
		if attempt.Before(since) {
			continue
		}

		sinceCount++
	}

	return sinceCount >= ch.attemptsMax
}
