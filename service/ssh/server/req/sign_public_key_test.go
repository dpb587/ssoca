package req_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	"github.com/dpb587/ssoca/certauth/memory/memoryfakes"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"
	serverreq "github.com/dpb587/ssoca/server/service/req"
	svcapi "github.com/dpb587/ssoca/service/ssh/api"
	svcconfig "github.com/dpb587/ssoca/service/ssh/server/config"
	. "github.com/dpb587/ssoca/service/ssh/server/req"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SignPublicKey", func() {
	var subject SignPublicKey
	var publicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDFkJcmlzl2PrJh0xv47AhQbiO2GSsUALlsN9EcA2kpN5dTqJOTiozxDWvJnwoIOLQqJuVbpSMC8BsgIekcJuI8CqQLCXmLL2CG84ltQiWnHNRAIcKv1Jqd7NUwwT9qmiGAX8mxDMLu0BOH7zXMcnRonLMVQz3G5tTK4jHIUiICFQ=="

	Describe("Route", func() {
		It("returns", func() {
			Expect(subject.Route()).To(Equal("sign-public-key"))
		})
	})

	Describe("Execute", func() {
		var fakecertauth certauthfakes.FakeProvider
		var realcertauth certauth.Provider
		var token *auth.Token
		var loggerContext logrus.Fields
		var res httptest.ResponseRecorder

		Context("common configuration", func() {
			BeforeEach(func() {
				loggerContext = logrus.Fields{
					"custom": "fake",
				}

				token = &auth.Token{ID: "fake-user"}
				fakecertauth = certauthfakes.FakeProvider{}
				realcertauth = memoryfakes.CreateMock1()
				res = *httptest.NewRecorder()

				subject = SignPublicKey{
					Validity: time.Duration(3600),
					Principals: dynamicvalue.MultiAnyValue{
						dynamicvalue.NewStringValue("vcap"),
					},
					CertAuth: &fakecertauth,
					CriticalOptions: svcconfig.CriticalOptions{
						svcconfig.CriticalOptionSourceAddress: "127.0.0.1",
					},
					Extensions: svcconfig.Extensions{
						svcconfig.ExtensionPermitAgentForwarding,
						svcconfig.ExtensionPermitPTY,
					},
					Target: svcconfig.Target{
						Host: "ssh.example.com",
					},
				}
			})

			It("works", func() {
				fakecertauth.SignSSHCertificateStub = realcertauth.SignSSHCertificate

				req := serverreq.Request{
					RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(fmt.Sprintf(`{"public_key":"%s"}`, strings.Replace(publicKey, "\n", "\\n", -1)))),
					RawResponse:   &res,
					AuthToken:     token,
					LoggerContext: loggerContext,
				}

				err := subject.Execute(req)

				Expect(err).ToNot(HaveOccurred())

				var resPayload svcapi.SignPublicKeyResponse

				err = json.Unmarshal(res.Body.Bytes(), &resPayload)
				Expect(err).ToNot(HaveOccurred())

				// @todo improve?
				resSplit := strings.Split(resPayload.Certificate, " ")
				Expect(resSplit).To(HaveLen(2))
				Expect(resSplit[0]).To(Equal("ssh-rsa-cert-v01@openssh.com"))
				Expect(len(resSplit[1])).To(BeNumerically(">", 512))

				Expect(resPayload.Target).ToNot(BeNil())
				Expect(resPayload.Target.Host).To(Equal("ssh.example.com"))
				Expect(resPayload.Target.User).To(Equal(""))
				Expect(resPayload.Target.Port).To(Equal(0))

				Expect(fakecertauth.SignSSHCertificateCallCount()).To(Equal(1))

				cert, innerLoggerContext := fakecertauth.SignSSHCertificateArgsForCall(0)
				Expect(cert.KeyId).To(Equal("fake-user"))
				Expect(cert.CertType).To(Equal(uint32(ssh.UserCert)))
				Expect(cert.ValidPrincipals).To(Equal([]string{"vcap"}))
				Expect(cert.Permissions.CriticalOptions).To(HaveLen(1))
				Expect(cert.Permissions.CriticalOptions[string(svcconfig.CriticalOptionSourceAddress)]).To(Equal("127.0.0.1"))
				Expect(cert.Permissions.Extensions).To(HaveKey(string(svcconfig.ExtensionPermitAgentForwarding)))
				Expect(cert.Permissions.Extensions[string(svcconfig.ExtensionPermitAgentForwarding)]).To(Equal(""))
				Expect(cert.Permissions.Extensions).To(HaveKey(string(svcconfig.ExtensionPermitPTY)))
				Expect(cert.Permissions.Extensions[string(svcconfig.ExtensionPermitPTY)]).To(Equal(""))
				Expect(innerLoggerContext["custom"]).To(Equal("fake"))
			})

			Context("invalid public keys", func() {
				Context("invalid format", func() {
					It("errors", func() {
						req := serverreq.Request{
							RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(`{"public_key":"invalid"}`)),
							RawResponse:   &res,
							AuthToken:     token,
							LoggerContext: loggerContext,
						}

						err := subject.Execute(req)

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Invalid public key format"))

						apiError, ok := err.(apierr.Error)
						Expect(ok).To(BeTrue())

						Expect(apiError.Status).To(Equal(400))
						Expect(apiError.PublicError).To(Equal("Failed to read public key"))
					})
				})

				Context("invalid data", func() {
					It("errors", func() {
						req := serverreq.Request{
							RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(`{"public_key":"ssh-rsa ="}`)),
							RawResponse:   &res,
							AuthToken:     token,
							LoggerContext: loggerContext,
						}

						err := subject.Execute(req)

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Decoding public key"))

						apiError, ok := err.(apierr.Error)
						Expect(ok).To(BeTrue())

						Expect(apiError.Status).To(Equal(400))
						Expect(apiError.PublicError).To(Equal("Failed to decode public key"))
					})
				})

				Context("invalid ssh key", func() {
					It("errors", func() {
						req := serverreq.Request{
							RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(`{"public_key":"ssh-rsa data"}`)),
							RawResponse:   &res,
							AuthToken:     token,
							LoggerContext: loggerContext,
						}

						err := subject.Execute(req)

						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("Parsing public key"))

						apiError, ok := err.(apierr.Error)
						Expect(ok).To(BeTrue())

						Expect(apiError.Status).To(Equal(400))
						Expect(apiError.PublicError).To(Equal("Failed to parse public key"))
					})
				})
			})

			Context("certauth failure", func() {
				It("errors", func() {
					fakecertauth.SignSSHCertificateReturns(errors.New("fake-err"))

					req := serverreq.Request{
						RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(fmt.Sprintf(`{"public_key":"%s"}`, strings.Replace(publicKey, "\n", "\\n", -1)))),
						RawResponse:   &res,
						AuthToken:     token,
						LoggerContext: loggerContext,
					}

					err := subject.Execute(req)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-err"))
					Expect(err.Error()).To(ContainSubstring("Signing certificate"))
				})
			})
		})

		Context("dynamic values", func() {
			BeforeEach(func() {
				loggerContext = logrus.Fields{
					"custom": "fake",
				}

				token = &auth.Token{ID: "fake-user"}
				fakecertauth = certauthfakes.FakeProvider{}
				realcertauth = memoryfakes.CreateMock1()

				targetUser := dynamicvalue.ConfigValue{}
				targetUser.WithDefault(dynamicvalue.MustCreateTemplateValue("{{ .Token.ID }}-suffixed"))

				subject = SignPublicKey{
					Validity: time.Duration(3600),
					Principals: dynamicvalue.MultiAnyValue{
						dynamicvalue.NewStringValue("static"),
						dynamicvalue.MustCreateTemplateValue("{{ .Token.ID }}"),
						dynamicvalue.MustCreateTemplateValue("{{ if false }}something{{ end }}"),
					},
					CertAuth: &fakecertauth,
					CriticalOptions: svcconfig.CriticalOptions{
						svcconfig.CriticalOptionSourceAddress: "127.0.0.1",
					},
					Extensions: svcconfig.Extensions{
						svcconfig.ExtensionPermitAgentForwarding,
						svcconfig.ExtensionPermitPTY,
					},
					Target: svcconfig.Target{
						Host: "ssh.example.com",
						User: targetUser,
					},
				}
			})

			It("works", func() {
				fakecertauth.SignSSHCertificateStub = realcertauth.SignSSHCertificate

				req := serverreq.Request{
					RawRequest:    httptest.NewRequest("GET", "https://localhost/file?name=test1", strings.NewReader(fmt.Sprintf(`{"public_key":"%s"}`, strings.Replace(publicKey, "\n", "\\n", -1)))),
					RawResponse:   &res,
					AuthToken:     token,
					LoggerContext: loggerContext,
				}

				err := subject.Execute(req)

				Expect(err).ToNot(HaveOccurred())

				var resPayload svcapi.SignPublicKeyResponse

				err = json.Unmarshal(res.Body.Bytes(), &resPayload)
				Expect(err).ToNot(HaveOccurred())

				cert, _ := fakecertauth.SignSSHCertificateArgsForCall(0)
				Expect(cert.ValidPrincipals).To(HaveLen(2))
				Expect(cert.ValidPrincipals).To(ContainElement("static"))
				Expect(cert.ValidPrincipals).To(ContainElement("fake-user"))
				Expect(resPayload.Target).ToNot(BeNil())
				Expect(resPayload.Target.User).To(Equal("fake-user-suffixed"))
			})
		})
	})
})
