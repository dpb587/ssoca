package req

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/dpb587/ssoca/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/server/service/req"
	"golang.org/x/oauth2"
)

type Initiate struct {
	Config oauth2.Config
}

var _ req.RouteHandler = Initiate{}

func (h Initiate) Route() string {
	return "initiate"
}

func (h Initiate) Execute(w http.ResponseWriter, r *http.Request) error {
	s := make([]byte, 32)
	rand.Read(s)

	state := base64.URLEncoding.EncodeToString(s)

	http.SetCookie(
		w,
		&http.Cookie{
			Name:  config.CookieStateName,
			Value: state,
		},
	)

	clientPort := r.FormValue("client_port")

	if clientPort != "" {
		http.SetCookie(
			w,
			&http.Cookie{
				Name:  config.CookieClientPortName,
				Value: clientPort,
			},
		)
	}

	url := h.Config.AuthCodeURL(state)

	http.Redirect(w, r, url, 302)

	return nil
}
