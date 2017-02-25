package req_test

import (
	"context"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	oauth2supportconfig "github.com/dpb587/ssoca/authn/support/oauth2/config"
	. "github.com/dpb587/ssoca/authn/support/oauth2/req"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockClaims map[string]interface{}

func (mockClaims) Valid() error {
	return nil
}

type mockTransport struct {
	rt func(req *http.Request) (resp *http.Response, err error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return t.rt(req)
}

var _ = Describe("Callback", func() {
	var subject Callback

	Describe("Route", func() {
		subject = Callback{}

		Expect(subject.Route()).To(Equal("callback"))
	})

	Describe("Execute", func() {
		var w httptest.ResponseRecorder
		var rt func(r *http.Request) (w *http.Response, err error)
		var privateKey *rsa.PrivateKey

		BeforeEach(func() {
			var err error
			privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALx3tjO/le0iinRqpdMa3/+4lbYuPQOsJROsDo9D4gN3rvnqOiuI
MIydz3VZ+ErHf5+LCQa4elzMK6mGpIVrewcCAwEAAQJBAJfHZPHZ6fkWpyBIPxF7
BEhiNBeKt1J80UM9fmA8UAlvd3m003ANZAIe6/GEYCTVp49jJ0e66Fwhd+LOsFN2
cPkCIQDKSnqE8QKg45S/QosRuCpM2EKwscp0855wdRAyQFYeWwIhAO6BrRRbAriQ
8tafNqIeKBaKBLZlen0ataAPiKAEqo3FAiAR1bMrmVwT9zycCC/ephAEqmRm06X3
3aqwW4HMDGQLVQIgUZL6rp6eJKA23l8gIXys+2CDUhsNNOLAwhjuAsT1zH0CIQC/
Y5S1jt/4jLs8Xfh6zPs15EuvQKpUGYKAIHofIp0T3Q==
-----END RSA PRIVATE KEY-----`))
			if err != nil {
				panic(err)
			}
		})

		parseClaims := func(token string) mockClaims {
			claims := mockClaims{}

			_, err := jwt.ParseWithClaims(
				token,
				&claims,
				func(_ *jwt.Token) (interface{}, error) {
					return &privateKey.PublicKey, nil
				},
			)

			if err != nil {
				panic(err)
			}

			return claims
		}

		BeforeEach(func() {
			subject = Callback{
				Origin: "fake-origin",
				UserProfileLoader: func(_ *http.Client) (oauth2supportconfig.UserProfile, error) {
					return oauth2supportconfig.UserProfile{
						Username: "fake-user",
						Attributes: map[string]string{
							"attr1": "value1",
							"attr2": "value2",
						},
						Scopes: []string{
							"scope1",
							"scope2",
						},
					}, nil
				},
				Config: oauth2.Config{
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					Endpoint: oauth2.Endpoint{
						TokenURL: "https://oauth.example.com/token",
					},
				},
				Context: context.WithValue(
					context.Background(),
					oauth2.HTTPClient,
					&http.Client{
						Transport: &mockTransport{
							rt: func(r *http.Request) (w *http.Response, err error) {
								return rt(r)
							},
						},
					},
				),
				JWT: oauth2supportconfig.JWT{
					PrivateKey: *privateKey,
				},
			}

			w = *httptest.NewRecorder()
		})

		Context("happy path", func() {
			var r *http.Request

			BeforeEach(func() {
				r = httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=state12345")

				rt = func(r *http.Request) (*http.Response, error) {
					Expect(r.FormValue("code")).To(Equal("myauthtoken"))

					return &http.Response{
						StatusCode: 200,
						Header:     http.Header{},
						Body:       ioutil.NopCloser(strings.NewReader(`{"access_token":"fake-access-token","token_type":"fake-token-type","refresh_token":"fake-refresh-token"}`)),
					}, nil
				}
			})

			It("signs tokens and removes state cookie", func() {
				err := subject.Execute(r, &w)

				Expect(err).ToNot(HaveOccurred())

				Expect(w.Header()["Content-Type"]).To(HaveLen(1))
				Expect(w.Header()["Content-Type"][0]).To(Equal("text/plain"))

				Expect(w.Header()["Set-Cookie"]).To(HaveLen(2))
				Expect(w.Header()["Set-Cookie"][0]).To(Equal("ssoca_oauth_state=; Max-Age=0"))
				Expect(w.Header()["Set-Cookie"][1]).To(MatchRegexp("^Authorization=[^;]+; Path=/$"))

				claims := parseClaims(w.Body.String())

				Expect(claims["username"]).To(Equal("fake-user"))
				Expect(claims["attributes"].(map[string]interface{})["attr1"]).To(Equal("value1"))
				Expect(claims["attributes"].(map[string]interface{})["attr2"]).To(Equal("value2"))
				Expect(claims["scope"].([]interface{})).To(ContainElement("scope1"))
				Expect(claims["scope"].([]interface{})).To(ContainElement("scope2"))

				Expect(claims["aud"]).To(Equal("fake-origin"))
				Expect(claims["iss"]).To(Equal("fake-origin"))
				// @todo test other fields; exp nbf sub iss aud
			})

			Context("client port", func() {
				It("sends html", func() {
					r.Header.Add("Cookie", "ssoca_oauth_clientport=12345")

					err := subject.Execute(r, &w)

					Expect(err).ToNot(HaveOccurred())

					Expect(w.Header()["Content-Type"]).To(HaveLen(1))
					Expect(w.Header()["Content-Type"][0]).To(Equal("text/html"))

					Expect(w.Header()["Set-Cookie"]).To(HaveLen(3))
					Expect(w.Header()["Set-Cookie"][1]).To(MatchRegexp("^Authorization=[^;]+; Path=/$"))
					Expect(w.Header()["Set-Cookie"][2]).To(Equal("ssoca_oauth_clientport=; Max-Age=0"))

					body := w.Body.String()

					Expect(body).To(ContainSubstring(` action="http://127.0.0.1:12345"`))
					Expect(body).To(ContainSubstring(` value="/ui/auth-success.html"`))

					re := regexp.MustCompile(` name="token" value="([^"]+)"`)
					tokenMatch := re.FindStringSubmatch(body)

					claims := parseClaims(tokenMatch[1])

					Expect(claims["username"]).To(Equal("fake-user"))
				})
			})
		})

		Context("state cookie", func() {
			It("errors when missing", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)

				err := subject.Execute(r, &w)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("state cookie"))
			})

			It("errors when missing", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=mismatch12345")

				err := subject.Execute(r, &w)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("value does not match"))
			})
		})

		Context("exchanging token", func() {
			It("errors when upstream error", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=state12345")

				rt = func(r *http.Request) (*http.Response, error) {
					return nil, errors.New("fake-err")
				}

				err := subject.Execute(r, &w)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Exchanging token"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})

			It("errors when invalid upstream data", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=state12345")

				rt = func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Header:     http.Header{},
						Body:       ioutil.NopCloser(strings.NewReader(`{"access_token":"","token_type":"fake-token-type","refresh_token":"fake-refresh-token"}`)),
					}, nil
				}

				err := subject.Execute(r, &w)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid token"))
			})
		})

		Context("profile loader errors", func() {
			It("errors", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=state12345")

				rt = func(r *http.Request) (*http.Response, error) {
					Expect(r.FormValue("code")).To(Equal("myauthtoken"))

					return &http.Response{
						StatusCode: 200,
						Header:     http.Header{},
						Body:       ioutil.NopCloser(strings.NewReader(`{"access_token":"fake-access-token","token_type":"fake-token-type","refresh_token":"fake-refresh-token"}`)),
					}, nil
				}

				subject.UserProfileLoader = func(_ *http.Client) (oauth2supportconfig.UserProfile, error) {
					return oauth2supportconfig.UserProfile{}, errors.New("fake-err")
				}

				err := subject.Execute(r, &w)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Loading user profile"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})
	})
})
