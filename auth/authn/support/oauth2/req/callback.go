package req

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dpb587/ssoca/auth/authn/support/oauth2/config"
	"github.com/dpb587/ssoca/auth/authn/support/selfsignedjwt"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/oauth2"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

var clientRedirectTemplateV2 = template.Must(template.New("html").Parse(`
	<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>ssoca</title>
			<script>
				setTimeout(function () {document.getElementById("token").style.display = "block"}, 5000);
				window.addEventListener(
					"message",
					function (event)
					{
					  if (event.origin != "http://127.0.0.1:{{.Port}}") {
					    return;
						} else if (event.data.action != "done") {
							return;
						}

						self.location = "{{.Redirect}}";
					},
					false
				);
			</script>
			<style type="text/css">
				body {
					color: #666666;
					font-family: sans-serif;
				}
				iframe {
					border: none;
					height: 0;
					width: 0;
				}
				pre {
					color: #333333;
					display: none;
					overflow-wrap: break-word;
					white-space: pre-wrap;
					width: 100%;
					word-break: break-all;
				}
			</style>
			<noscript>
				<style type="text/css">
					pre {
						display: block !important;
					}
				</style>
			</noscript>
		</head>
		<body>
			<p>Sending token to ssoca at 127.0.0.1:{{.Port}}&hellip;</p>
			<pre id="token"><code>{{.Token}}</code></pre>
			<iframe src="http://127.0.0.1:{{.Port}}/" onload="this.contentWindow.postMessage({'action':'set_token','params':{'token':document.getElementById('token').innerText}}, 'http://127.0.0.1:{{.Port}}')"></iframe>
		</body>
	</html>
`))

// Historically we relied on POSTing (initially to avoid the chance of an
// authentication token ending up in browser history; also POST is more correct
// for this behavior); but some browsers (e.g. Safari) are giving errors about
// redirecting to an insecure form. To avoid relying on JavaScript, this moved
// to using GET and postMessage.
var clientRedirectTemplatePOST = template.Must(template.New("html").Parse(`
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
		return apierr.NewError(apierr.WrapError(err, "Getting state cookie"), http.StatusBadRequest, "State cookie does not exist")
	}

	if request.RawRequest.URL.Query().Get("state") != state.Value {
		return apierr.NewError(errors.New("State cookie value does not match expected state"), http.StatusBadRequest, "State cookie does not match")
	}

	oauthToken, err := h.Config.Exchange(h.Context, request.RawRequest.URL.Query().Get("code"))
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
		return bosherr.WrapError(err, "Signing token")
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

		var clientRedirectTemplate *template.Template = clientRedirectTemplateV2

		clientVersion, _ := request.RawRequest.Cookie(config.CookieClientVersionName)

		if clientVersion != nil {
			switch clientVersion.Value {
			case "0.13.0", "0.12.0", "0.11.0", "0.10.0", "0.9.0", "0.8.0", "0.7.1", "0.7.0":
				// hard-coding these older versions as simpler-lazier than version parsing
				clientRedirectTemplate = clientRedirectTemplatePOST
			}
		}

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
