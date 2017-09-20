package req

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
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
	redirect_uri, err := url.Parse(h.Config.RedirectURL)
	if err != nil {
		return bosherr.WrapError(err, "Parsing redirect URL")
	}

	if req.RawRequest.Host != redirect_uri.Host {
		target := req.RawRequest.URL
		target.Host = redirect_uri.Host
		target.Scheme = redirect_uri.Scheme

		http.Redirect(req.RawResponse, req.RawRequest, target.String(), 302)

		return nil
	}

	s := make([]byte, 32)
	rand.Read(s)

	state := base64.URLEncoding.EncodeToString(s)

	http.SetCookie(
		req.RawResponse,
		&http.Cookie{
			Domain: redirect_uri.Hostname(),
			Name:   config.CookieStateName,
			Path:   "/auth/",
			Secure: redirect_uri.Scheme == "https",
			Value:  state,
		},
	)

	clientPort := req.RawRequest.FormValue("client_port")

	if clientPort != "" {
		http.SetCookie(
			req.RawResponse,
			&http.Cookie{
				Domain: redirect_uri.Hostname(),
				Name:   config.CookieClientPortName,
				Path:   "/auth/",
				Secure: redirect_uri.Scheme == "https",
				Value:  clientPort,
			},
		)
	}

	url := h.Config.AuthCodeURL(state)

	http.Redirect(req.RawResponse, req.RawRequest, url, 302)

	return nil
}
