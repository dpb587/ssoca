package req

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
)

var clientRedirectTemplate = template.Must(template.New("html").Parse(`
	<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>ssoca</title>
		</head>
		<body onload="document.getElementsByTagName('input')[0].click();">
			<noscript>
				<pre><code>{{.Token}}</code></pre>
			</noscript>
			<form method="post" action="http://127.0.0.1:{{.Port}}">
				<input type="submit" style="display:none;" />
				<input type="hidden" name="token" value="{{.Token}}" />
				<input type="hidden" name="return_to" value="{{.Redirect}}" />
			</form>
		</body>
	</html>
`))

type Callback struct {
	URLs              config.URLs
	UserProfileLoader config.UserProfileLoader
	Config            oauth2.Config
	Context           context.Context
	JWT               config.JWT

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Callback{}

func (h Callback) Route() string {
	return "callback"
}

func (h Callback) Execute(request req.Request) error {
	state, err := request.RawRequest.Cookie(config.CookieStateName)
	if err != nil {
		return apierr.NewError(apierr.WrapError(err, "getting state cookie"), http.StatusBadRequest, "state cookie does not exist")
	}

	if request.RawRequest.URL.Query().Get("state") != state.Value {
		return apierr.NewError(errors.New("state cookie value does not match expected state"), http.StatusBadRequest, "state cookie does not match")
	}

	oauthToken, err := h.Config.Exchange(h.Context, request.RawRequest.URL.Query().Get("code"))
	if err != nil {
		return errors.Wrap(err, "exchanging token")
	}

	if !oauthToken.Valid() {
		return errors.New("invalid token")
	}

	userProfile, err := h.UserProfileLoader(h.Config.Client(h.Context, oauthToken))
	if err != nil {
		return errors.Wrap(err, "loading user profile")
	}

	tokenUUID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "generating local token ID")
	}

	token := jwt.NewWithClaims(config.JWTSigningMethod, selfsignedjwt.Token{
		ID:         userProfile.ID,
		Groups:     userProfile.Groups,
		Attributes: userProfile.Attributes,
		StandardClaims: jwt.StandardClaims{
			Audience:  h.URLs.Origin,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        tokenUUID.String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    h.URLs.Origin,
			NotBefore: time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString(&h.JWT.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "signing token")
	}

	// remove the cookie
	http.SetCookie(
		request.RawResponse,
		&http.Cookie{
			Name:   state.Name,
			MaxAge: -1,
		},
	)

	clientPort, _ := request.RawRequest.Cookie(config.CookieClientPortName)

	if clientPort != nil {
		http.SetCookie(
			request.RawResponse,
			&http.Cookie{
				Name:   clientPort.Name,
				MaxAge: -1,
			},
		)

		var buf bytes.Buffer

		err = clientRedirectTemplate.Execute(
			&buf,
			struct {
				Token    string
				Port     string
				Redirect string
			}{
				Token:    tokenString,
				Port:     clientPort.Value,
				Redirect: h.URLs.AuthSuccess,
			},
		)
		if err != nil {
			return err
		}

		request.RawResponse.Header().Set("Content-Type", "text/html")
		request.RawResponse.Write(buf.Bytes())
	} else {
		request.RawResponse.Header().Set("Content-Type", "text/plain")
		request.RawResponse.Write([]byte(tokenString))
	}

	return nil
}
