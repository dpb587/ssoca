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
	"github.com/dpb587/ssoca/auth"
	internaltests "github.com/dpb587/ssoca/auth/authn/oauth2/internal/tests"
	oauth2supportconfig "github.com/dpb587/ssoca/auth/authn/oauth2/server/config"
	. "github.com/dpb587/ssoca/auth/authn/oauth2/server/req"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
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
		It("works", func() {
			subject = Callback{}

			Expect(subject.Route()).To(Equal("callback"))
		})
	})

	Describe("Execute", func() {
		var res httptest.ResponseRecorder
		var rt func(r *http.Request) (w *http.Response, err error)
		var privateKey *rsa.PrivateKey

		BeforeEach(func() {
			var err error
			privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(internaltests.SharedPrivateKey))
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
			user := "fake-user"
			subject = Callback{
				URLs: oauth2supportconfig.URLs{
					Origin:      "fake-origin",
					AuthSuccess: "/fake-ui/auth-success.html",
				},
				UserProfileLoader: func(_ *http.Client) (auth.Token, error) {
					return auth.Token{
						ID: "fake-id",
						Groups: []string{
							"scope1",
							"scope2",
						},
						Attributes: map[auth.TokenAttribute]*string{
							auth.TokenUsernameAttribute: &user,
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
					PrivateKey: &oauth2supportconfig.PrivateKey{
						PrivateKey: privateKey,
					},
				},
			}

			res = *httptest.NewRecorder()
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
				err := subject.Execute(req.Request{
					RawRequest:  r,
					RawResponse: &res,
				})

				Expect(err).ToNot(HaveOccurred())

				Expect(res.Header()["Content-Type"]).To(HaveLen(1))
				Expect(res.Header()["Content-Type"][0]).To(Equal("text/plain"))

				Expect(res.Header()["Set-Cookie"]).To(HaveLen(1))
				Expect(res.Header()["Set-Cookie"][0]).To(Equal("ssoca_oauth_state=; Max-Age=0"))

				claims := parseClaims(res.Body.String())

				Expect(claims["scid"]).To(Equal("fake-id"))
				Expect(claims["scat"]).To(Equal(map[string]interface{}{"username": "fake-user"}))
				Expect(claims["scgr"].([]interface{})).To(ContainElement("scope1"))
				Expect(claims["scgr"].([]interface{})).To(ContainElement("scope2"))

				Expect(claims["aud"]).To(Equal("fake-origin"))
				Expect(claims["iss"]).To(Equal("fake-origin"))
				// @todo test other fields; exp nbf sub iss aud
			})

			Context("client port", func() {
				It("sends html", func() {
					r.Header.Add("Cookie", "ssoca_oauth_clientport=12345")

					err := subject.Execute(req.Request{
						RawRequest:  r,
						RawResponse: &res,
					})

					Expect(err).ToNot(HaveOccurred())

					Expect(res.Header()["Content-Type"]).To(HaveLen(1))
					Expect(res.Header()["Content-Type"][0]).To(Equal("text/html"))

					Expect(res.Header()["Set-Cookie"]).To(HaveLen(2))
					Expect(res.Header()["Set-Cookie"][0]).To(Equal("ssoca_oauth_state=; Max-Age=0"))
					Expect(res.Header()["Set-Cookie"][1]).To(Equal("ssoca_oauth_clientport=; Max-Age=0"))

					body := res.Body.String()

					Expect(body).To(ContainSubstring(` action="http://127.0.0.1:12345"`))
					Expect(body).To(ContainSubstring(` value="/fake-ui/auth-success.html"`))

					re := regexp.MustCompile(` name="token" value="([^"]+)"`)
					tokenMatch := re.FindStringSubmatch(body)

					claims := parseClaims(tokenMatch[1])

					Expect(claims["scid"]).To(Equal("fake-id"))
				})
			})
		})

		Context("state cookie", func() {
			It("errors when missing", func() {
				err := subject.Execute(req.Request{
					RawRequest:  httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil),
					RawResponse: &res,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("state cookie"))

				errapi, ok := err.(apierr.Error)
				Expect(ok).To(BeTrue())
				Expect(errapi.Status).To(Equal(400))
				Expect(errapi.PublicError).To(Equal("state cookie does not exist"))
			})

			It("errors when missing", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=mismatch12345")

				err := subject.Execute(req.Request{
					RawRequest:  r,
					RawResponse: &res,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("value does not match"))

				errapi, ok := err.(apierr.Error)
				Expect(ok).To(BeTrue())
				Expect(errapi.Status).To(Equal(400))
				Expect(errapi.PublicError).To(Equal("state cookie does not match"))
			})
		})

		Context("exchanging token", func() {
			It("errors when upstream error", func() {
				r := httptest.NewRequest("GET", "http://localhost/auth/callback?state=state12345&code=myauthtoken", nil)
				r.Header.Add("Cookie", "ssoca_oauth_state=state12345")

				rt = func(r *http.Request) (*http.Response, error) {
					return nil, errors.New("fake-err")
				}

				err := subject.Execute(req.Request{
					RawRequest:  r,
					RawResponse: &res,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("exchanging token"))
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

				err := subject.Execute(req.Request{
					RawRequest:  r,
					RawResponse: &res,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid token"))
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

				subject.UserProfileLoader = func(_ *http.Client) (auth.Token, error) {
					return auth.Token{}, errors.New("fake-err")
				}

				err := subject.Execute(req.Request{
					RawRequest:  r,
					RawResponse: &res,
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("loading user profile"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})
	})
})
