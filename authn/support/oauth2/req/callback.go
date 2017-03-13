package req

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/authn/support/selfsignedjwt"
	"github.com/dpb587/ssoca/server/service/req"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/oauth2"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
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
	Origin            string
	UserProfileLoader config.UserProfileLoader
	Config            oauth2.Config
	Context           context.Context
	JWT               config.JWT
}

var _ req.RouteHandler = Callback{}

func (h Callback) Route() string {
	return "callback"
}

func (h Callback) Execute(r *http.Request, w http.ResponseWriter) error {
	state, err := r.Cookie(config.CookieStateName)
	if err != nil {
		return bosherr.WrapError(err, "Getting state cookie")
	}

	if r.URL.Query().Get("state") != state.Value {
		return errors.New("State cookie value does not match expected state")
	}

	oauthToken, err := h.Config.Exchange(h.Context, r.URL.Query().Get("code"))
	if err != nil {
		return bosherr.WrapError(err, "Exchanging token")
	}

	if !oauthToken.Valid() {
		return errors.New("Invalid token")
	}

	userProfile, err := h.UserProfileLoader(h.Config.Client(h.Context, oauthToken))
	if err != nil {
		return bosherr.WrapError(err, "Loading user profile")
	}

	tokenUUID, err := uuid.NewV4()
	if err != nil {
		return bosherr.WrapError(err, "Generating local token ID")
	}

	token := jwt.NewWithClaims(config.JWTSigningMethod, selfsignedjwt.Token{
		ID:         userProfile.ID,
		Groups:     userProfile.Groups,
		Attributes: userProfile.Attributes,
		StandardClaims: jwt.StandardClaims{
			Audience:  h.Origin,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        tokenUUID.String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    h.Origin,
			NotBefore: time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString(&h.JWT.PrivateKey)
	if err != nil {
		return bosherr.WrapError(err, "Signing token")
	}

	// remove the cookie
	http.SetCookie(
		w,
		&http.Cookie{
			Name:   state.Name,
			MaxAge: -1,
		},
	)

	// ui cookie
	// @todo configurable
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "Authorization",
			Value: fmt.Sprintf("bearer %s", tokenString),
			Path:  "/",
		},
	)

	clientPort, _ := r.Cookie(config.CookieClientPortName)

	if clientPort != nil {
		http.SetCookie(
			w,
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
				Redirect: "/ui/auth-success.html", // @todo
			},
		)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(buf.Bytes())
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(tokenString))
	}

	return nil
}
