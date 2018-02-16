package req

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/crewjam/saml"
	svcconfig "github.com/dpb587/ssoca/auth/authn/saml/config"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/server/service/req"
)

type Initiate struct {
	Config svcconfig.Config
	IDP    *saml.ServiceProvider

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Initiate{}

func (h Initiate) Route() string {
	return "initiate"
}

func (h Initiate) Execute(req req.Request) error {
	redirect_uri := h.IDP.AcsURL

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

	clientVersion := req.RawRequest.FormValue("client_version")

	if clientVersion != "" {
		http.SetCookie(
			req.RawResponse,
			&http.Cookie{
				Domain: redirect_uri.Hostname(),
				Name:   config.CookieClientVersionName,
				Path:   "/",
				Secure: redirect_uri.Scheme == "https",
				Value:  clientVersion,
			},
		)
	}

	authreq, err := h.IDP.MakePostAuthenticationRequest("TODO")
	if err != nil {
		return bosherr.WrapError(err, "Creating authentication request")
	}

	_, err = req.RawResponse.Write(authreq)
	if err != nil {
		return bosherr.WrapError(err, "Writing authentication request")
	}

	return nil
}
