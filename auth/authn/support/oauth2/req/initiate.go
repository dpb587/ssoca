package req

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/server/service/req"
	"golang.org/x/oauth2"
)

type Initiate struct {
	Config oauth2.Config

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Initiate{}

func (h Initiate) Route() string {
	return "initiate"
}

func (h Initiate) Execute(req req.Request) error {
	s := make([]byte, 32)
	rand.Read(s)

	state := base64.URLEncoding.EncodeToString(s)

	http.SetCookie(
		req.RawResponse,
		&http.Cookie{
			Name:  config.CookieStateName,
			Value: state,
		},
	)

	clientPort := req.RawRequest.FormValue("client_port")

	if clientPort != "" {
		http.SetCookie(
			req.RawResponse,
			&http.Cookie{
				Name:  config.CookieClientPortName,
				Value: clientPort,
			},
		)
	}

	url := h.Config.AuthCodeURL(state)

	http.Redirect(req.RawResponse, req.RawRequest, url, 302)

	return nil
}
