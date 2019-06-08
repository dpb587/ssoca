package server_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/dpb587/ssoca/service/githubauth/server"
	svcconfig "github.com/dpb587/ssoca/service/githubauth/server/config"
	"golang.org/x/oauth2"

	oauth2server "github.com/dpb587/ssoca/auth/authn/support/oauth2/server"
	oauth2config "github.com/dpb587/ssoca/auth/authn/support/oauth2/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockTransport struct {
	rt func(req *http.Request) (resp *http.Response, err error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return t.rt(req)
}

var _ = Describe("Auth", func() {
	var subject *Service

	Describe("OAuthUserProfileLoader", func() {
		Context("simple config", func() {
			BeforeEach(func() {
				subject = NewService(
					"auth",
					svcconfig.Config{},
					oauth2server.NewService(
						oauth2config.URLs{Origin: "test"},
						oauth2.Config{},
						oauth2.NoContext,
						oauth2config.JWT{},
					),
				)
			})

			It("works", func() {
				profile, err := subject.OAuthUserProfileLoader(&http.Client{
					Transport: &mockTransport{
						rt: func(r *http.Request) (w *http.Response, err error) {
							switch r.URL.String() {
							case "https://api.github.com/user":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"login":"octocat","name":"monalisa octocat"}`)),
								}, nil
							case "https://api.github.com/user/teams?page=1":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`[{"slug":"demo-one","organization":{"login":"test1"}},{"slug":"demo-two","organization":{"login":"test1"}},{"slug":"demo-three","organization":{"login":"test2"}}]`)),
								}, nil
							case "https://api.github.com/user/orgs?page=1":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`[{"login":"test1"},{"login":"test2"},{"login":"test3"}]`)),
								}, nil
							}

							Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

							return &http.Response{}, nil
						},
					},
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(profile.Username()).To(Equal("octocat"))
				Expect(profile.Groups).To(HaveLen(6))
				Expect(profile.Groups).To(ContainElement("test1/demo-one"))
				Expect(profile.Groups).To(ContainElement("test1/demo-two"))
				Expect(profile.Groups).To(ContainElement("test2/demo-three"))
				Expect(profile.Groups).To(ContainElement("test1"))
				Expect(profile.Groups).To(ContainElement("test2"))
				Expect(profile.Groups).To(ContainElement("test3"))
				Expect(profile.Name()).To(Equal("monalisa octocat"))
			})

			Context("bad user info requests", func() {
				Context("transport errors", func() {
					It("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									return nil, errors.New("fake-err")
								},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("fetching user info"))
						Expect(err.Error()).To(ContainSubstring("fake-err"))
					})
				})

				Context("server errors", func() {
					XIt("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									switch r.URL.String() {
									case "https://api.github.com/user":
										return &http.Response{
											StatusCode: 400,
											Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
										}, nil
									}

									Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

									return &http.Response{}, nil
								},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("failed to request user info"))
					})
				})
			})
		})
	})
})
