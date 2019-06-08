package server_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/dpb587/ssoca/service/googleauth/server"
	svcconfig "github.com/dpb587/ssoca/service/googleauth/server/config"

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
					"http://example.com",
					svcconfig.Config{},
				)
			})

			It("works", func() {
				profile, err := subject.OAuthUserProfileLoader(&http.Client{
					Transport: &mockTransport{
						rt: func(r *http.Request) (w *http.Response, err error) {
							switch r.URL.String() {
							case "https://www.googleapis.com/oauth2/v3/userinfo":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
								}, nil
							}

							Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

							return &http.Response{}, nil
						},
					},
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(profile.Username()).To(Equal("somebody@example.com"))
				Expect(profile.Groups).To(HaveLen(3))
				Expect(profile.Groups).To(ContainElement("somebody@example.com"))
				Expect(profile.Groups).To(ContainElement("email/mailbox/somebody"))
				Expect(profile.Groups).To(ContainElement("email/domain/example.com"))
				Expect(profile.Name()).To(Equal("test user"))
			})

			Context("with non-verified email", func() {
				It("errors", func() {
					_, err := subject.OAuthUserProfileLoader(&http.Client{
						Transport: &mockTransport{
							rt: func(r *http.Request) (w *http.Response, err error) {
								switch r.URL.String() {
								case "https://www.googleapis.com/oauth2/v3/userinfo":
									return &http.Response{
										StatusCode: 200,
										Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":false}`)),
									}, nil
								}

								Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

								return &http.Response{}, nil
							},
						},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unverified email"))
				})
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
					It("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									switch r.URL.String() {
									case "https://www.googleapis.com/oauth2/v3/userinfo":
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

				Context("bad server response data", func() {
					It("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									switch r.URL.String() {
									case "https://www.googleapis.com/oauth2/v3/userinfo":
										return &http.Response{
											StatusCode: 200,
											Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody`)),
										}, nil
									}

									Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

									return &http.Response{}, nil
								},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unmarshaling user info"))
					})
				})
			})
		})

		Context("extended cloud project config", func() {
			var cloudProjectConfig svcconfig.ScopesCloudProjectConfig

			BeforeEach(func() {
				cloudProjectConfig = svcconfig.ScopesCloudProjectConfig{}

				subject = NewService(
					"auth",
					"http://example.com",
					svcconfig.Config{
						Scopes: svcconfig.ScopesConfig{
							CloudProject: &cloudProjectConfig,
						},
					},
				)
			})

			It("works", func() {
				profile, err := subject.OAuthUserProfileLoader(&http.Client{
					Transport: &mockTransport{
						rt: func(r *http.Request) (w *http.Response, err error) {
							switch r.URL.String() {
							case "https://www.googleapis.com/oauth2/v3/userinfo":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
								}, nil
							case "https://cloudresourcemanager.googleapis.com/v1/projects?alt=json&pageSize=1024":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"projects":[{"projectId":"test-1234","projectNumber":"12345678","name":"Test 1234"},{"projectId":"test-2345","projectNumber":"23456789","name":"Test 2345"}]}`)),
								}, nil
							case "https://cloudresourcemanager.googleapis.com/v1/projects/test-1234:getIamPolicy?alt=json":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"bindings":[{"role":"roles/owner","members":["user:somebody@example.com"]}]}`)),
								}, nil
							case "https://cloudresourcemanager.googleapis.com/v1/projects/test-2345:getIamPolicy?alt=json":
								return &http.Response{
									StatusCode: 200,
									Body:       ioutil.NopCloser(strings.NewReader(`{"bindings":[{"role":"roles/editor","members":["user:somebody@example.com"]}]}`)),
								}, nil
							}

							Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

							return &http.Response{}, nil
						},
					},
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(profile.Username()).To(Equal("somebody@example.com"))
				Expect(profile.Groups).To(HaveLen(5))
				Expect(profile.Groups).To(ContainElement("somebody@example.com"))
				Expect(profile.Groups).To(ContainElement("email/mailbox/somebody"))
				Expect(profile.Groups).To(ContainElement("email/domain/example.com"))
				Expect(profile.Groups).To(ContainElement("cloud/project/test-1234/roles/owner"))
				Expect(profile.Groups).To(ContainElement("cloud/project/test-2345/roles/editor"))
				Expect(profile.Name()).To(Equal("test user"))
			})

			Context("filtering project/role", func() {
				var client *http.Client

				BeforeEach(func() {
					client = &http.Client{
						Transport: &mockTransport{
							rt: func(r *http.Request) (w *http.Response, err error) {
								switch r.URL.String() {
								case "https://www.googleapis.com/oauth2/v3/userinfo":
									return &http.Response{
										StatusCode: 200,
										Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
									}, nil
								case "https://cloudresourcemanager.googleapis.com/v1/projects?alt=json&pageSize=1024":
									return &http.Response{
										StatusCode: 200,
										Body:       ioutil.NopCloser(strings.NewReader(`{"projects":[{"projectId":"test-1234","projectNumber":"12345678","name":"Test 1234"},{"projectId":"test-2345","projectNumber":"23456789","name":"Test 2345"}]}`)),
									}, nil
								case "https://cloudresourcemanager.googleapis.com/v1/projects/test-1234:getIamPolicy?alt=json":
									return &http.Response{
										StatusCode: 200,
										Body:       ioutil.NopCloser(strings.NewReader(`{"bindings":[{"role":"roles/owner","members":["user:somebody@example.com"]}]}`)),
									}, nil
								case "https://cloudresourcemanager.googleapis.com/v1/projects/test-2345:getIamPolicy?alt=json":
									return &http.Response{
										StatusCode: 200,
										Body:       ioutil.NopCloser(strings.NewReader(`{"bindings":[{"role":"roles/editor","members":["user:somebody@example.com"]}]}`)),
									}, nil
								}

								Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

								return &http.Response{}, nil
							},
						},
					}
				})

				Context("filtering project", func() {
					BeforeEach(func() {
						cloudProjectConfig.Projects = []string{"test-1234"}
					})

					It("works", func() {
						profile, err := subject.OAuthUserProfileLoader(client)

						Expect(err).ToNot(HaveOccurred())
						Expect(profile.Username()).To(Equal("somebody@example.com"))
						Expect(profile.Groups).To(HaveLen(4))
						Expect(profile.Groups).To(ContainElement("somebody@example.com"))
						Expect(profile.Groups).To(ContainElement("email/mailbox/somebody"))
						Expect(profile.Groups).To(ContainElement("email/domain/example.com"))
						Expect(profile.Groups).To(ContainElement("cloud/project/test-1234/roles/owner"))
						Expect(profile.Name()).To(Equal("test user"))
					})
				})

				Context("filtering role", func() {
					BeforeEach(func() {
						cloudProjectConfig.Roles = []string{"roles/editor"}
					})

					It("works", func() {
						profile, err := subject.OAuthUserProfileLoader(client)

						Expect(err).ToNot(HaveOccurred())
						Expect(profile.Username()).To(Equal("somebody@example.com"))
						Expect(profile.Groups).To(HaveLen(4))
						Expect(profile.Groups).To(ContainElement("somebody@example.com"))
						Expect(profile.Groups).To(ContainElement("email/mailbox/somebody"))
						Expect(profile.Groups).To(ContainElement("email/domain/example.com"))
						Expect(profile.Groups).To(ContainElement("cloud/project/test-2345/roles/editor"))
						Expect(profile.Name()).To(Equal("test user"))
					})
				})
			})

			Context("api request failures", func() {
				Context("project listing", func() {
					It("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									switch r.URL.String() {
									case "https://www.googleapis.com/oauth2/v3/userinfo":
										return &http.Response{
											StatusCode: 200,
											Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
										}, nil
									case "https://cloudresourcemanager.googleapis.com/v1/projects?alt=json&pageSize=1024":
										return &http.Response{
											StatusCode: 500,
											Body:       ioutil.NopCloser(strings.NewReader("error")),
										}, nil
									}

									Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

									return &http.Response{}, nil
								},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("listing projects"))
					})
				})

				Context("iam policy listing", func() {
					It("errors", func() {
						_, err := subject.OAuthUserProfileLoader(&http.Client{
							Transport: &mockTransport{
								rt: func(r *http.Request) (w *http.Response, err error) {
									switch r.URL.String() {
									case "https://www.googleapis.com/oauth2/v3/userinfo":
										return &http.Response{
											StatusCode: 200,
											Body:       ioutil.NopCloser(strings.NewReader(`{"name":"test user","email":"somebody@example.com","email_verified":true}`)),
										}, nil
									case "https://cloudresourcemanager.googleapis.com/v1/projects?alt=json&pageSize=1024":
										return &http.Response{
											StatusCode: 200,
											Body:       ioutil.NopCloser(strings.NewReader(`{"projects":[{"projectId":"test-1234","projectNumber":"12345678","name":"Test 1234"},{"projectId":"test-2345","projectNumber":"23456789","name":"Test 2345"}]}`)),
										}, nil
									case "https://cloudresourcemanager.googleapis.com/v1/projects/test-1234:getIamPolicy?alt=json":
										return &http.Response{
											StatusCode: 500,
											Body:       ioutil.NopCloser(strings.NewReader("error")),
										}, nil
									}

									Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

									return &http.Response{}, nil
								},
							},
						})

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("getting IAM policy"))
					})
				})
			})
		})
	})
})
