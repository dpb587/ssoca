package management

import (
	"encoding/base64"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/openvpn/client/profile"
)

const attemptsUntilRenewal = 3
const attemptsUntilRestart = 5

type needCertificateAction int

var needCertificateGetProfileAction needCertificateAction = 01
var needCertificateRenewalAction needCertificateAction = 010
var needCertificateRestartAction needCertificateAction = 0100

type DefaultHandler struct {
	profileManager profile.Manager

	attempts       []time.Time
	attemptsOffset int
	attemptsMax    int
}

var _ ServerHandler = &DefaultHandler{}

func NewDefaultHandler(profileManager profile.Manager) *DefaultHandler {
	attemptsMax := attemptsUntilRestart

	return &DefaultHandler{
		profileManager: profileManager,
		attempts:       make([]time.Time, attemptsMax),
		attemptsMax:    attemptsMax,
	}
}

func (ch *DefaultHandler) NeedCertificate(w io.Writer, _ string) (ServerHandlerCallback, error) {
	var err error

	// internally we want to throttle openvpn when it's doing repeated certificate
	// lookups
	action := ch.planNeedCertificate()

	if action&needCertificateRestartAction > 0 {
		// we tried our best, but something seems terribly wrong after quite a few
		// failures. tell openvpn to restart and hope that helps things.
		w.Write([]byte("signal SIGHUP\n"))

		return SuccessCallback, nil
	} else if action&needCertificateRenewalAction > 0 {
		// we failed a couple times, so let's optimistically assume openvn no longer
		// likes our certificate and explicitly get a new one. Note that GetProfile
		// already renews certificates when they expire, but the server might do
		// additional time-based checks on it to force shorter use validities.
		err := ch.profileManager.Renew()
		if err != nil {
			return nil, errors.Wrap(err, "renewing profile")
		}
	}

	var ovpn profile.Profile

	if action&needCertificateGetProfileAction > 0 {
		ovpn, err = ch.profileManager.GetProfile()
		if err != nil {
			return nil, errors.Wrap(err, "retrieving profile")
		}
	}

	w.Write([]byte("certificate\n"))
	w.Write(ovpn.CertificatePEM())
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
		return nil, errors.Wrap(err, "decoding signing token")
	}

	signature, err := ch.profileManager.Sign(data64)
	if err != nil {
		return nil, errors.Wrap(err, "signing token")
	}

	signature64 := base64.StdEncoding.EncodeToString(signature)

	w.Write([]byte("rsa-sig\n"))
	w.Write([]byte(signature64))
	w.Write([]byte("\n"))
	w.Write([]byte("END\n"))

	return SuccessCallback, nil
}

func (ch *DefaultHandler) planNeedCertificate() needCertificateAction {
	ch.attemptsOffset = (ch.attemptsOffset + 1) % ch.attemptsMax
	ch.attempts[ch.attemptsOffset] = time.Now()

	// fairly arbitrary time period: openvpn usually waits 5 seconds before
	// reconnection, and then allow some extra pre-failure connection negotiation
	// time.
	since := time.Now().Add(-3 * time.Minute)
	sinceCount := 0

	for _, attempt := range ch.attempts {
		if attempt.Before(since) {
			continue
		}

		sinceCount++
	}

	if sinceCount == attemptsUntilRenewal {
		return needCertificateRenewalAction | needCertificateGetProfileAction
	} else if sinceCount >= attemptsUntilRestart {
		return needCertificateRestartAction
	}

	return needCertificateGetProfileAction
}
