package client_test

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"github.com/dpb587/ssoca/auth/authn"
	. "github.com/dpb587/ssoca/auth/authn/support/oauth2/client"

	"github.com/dpb587/ssoca/client/clientfakes"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/config/configfakes"

	uifakes "github.com/cloudfoundry/bosh-cli/ui/fakes"
	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthService", func() {
	var subject *AuthService

	var env config.EnvironmentState
	var ui *uifakes.FakeUI
	var cmdRunner *boshsysfakes.FakeCmdRunner

	var runtime *clientfakes.FakeRuntime
	var configManager *configfakes.FakeManager

	BeforeEach(func() {
		env = config.EnvironmentState{
			Auth: &config.EnvironmentAuthState{},
			URL:  "https://192.0.2.99",
		}

		ui = &uifakes.FakeUI{}
		cmdRunner = &boshsysfakes.FakeCmdRunner{}

		configManager = &configfakes.FakeManager{}

		runtime = &clientfakes.FakeRuntime{}
		runtime.GetUIReturns(ui)
		runtime.GetStderrReturns(&bytes.Buffer{})
		runtime.GetStdoutReturns(&bytes.Buffer{})
		// runtime.GetStdinReturns(bytes.NewBufferString(""))
		runtime.GetConfigManagerReturns(configManager, nil)

		subject = NewAuthService("fake-name", "fake-type", runtime, cmdRunner)
	})

	Describe("AuthLogin", func() {
		// TODO enable tests; some complications around fake stdin/ui.AskedText
		XIt("shows login url and retrieves token from user", func() {
			runtime.GetEnvironmentReturns(env, nil)

			ui.AskedText = append(ui.AskedText, uifakes.Answer{
				Text:  "mytoken",
				Error: nil,
			})

			err := subject.AuthLogin()

			Expect(err).ToNot(HaveOccurred())
			//
			// authConfig, ok := token.(authn.AuthorizationToken)
			// Expect(ok).To(BeTrue())
			//
			// Expect(authConfig.Token).To(Equal("mytoken"))

			Expect(strings.Join(ui.Said, "")).To(ContainSubstring("https://192.0.2.99/auth/initiate"))
		})

		XIt("propagates errors on input fail", func() {
			runtime.GetEnvironmentReturns(env, nil)

			ui.AskedText = append(ui.AskedText, uifakes.Answer{
				Error: errors.New("ctrl c"),
			})

			err := subject.AuthLogin()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ctrl c"))
		})
	})

	Describe("AuthLogout", func() {
		It("does not error", func() {
			err := subject.AuthLogout()

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("AuthRequest", func() {
		var req http.Request

		BeforeEach(func() {
			req = http.Request{
				Header: http.Header{},
			}
		})

		It("adds deprecated authorization header", func() {
			env.Auth.Options = authn.AuthorizationToken{
				Token: "mytoken",
			}

			runtime.GetEnvironmentReturns(env, nil)

			err := subject.AuthRequest(&req)

			Expect(err).ToNot(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("bearer mytoken"))
		})

		It("adds authorization header", func() {
			env.Auth.Options = authn.AuthorizationToken{
				Type:  "Bearer",
				Value: "fake-value",
			}

			runtime.GetEnvironmentReturns(env, nil)

			err := subject.AuthRequest(&req)

			Expect(err).ToNot(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("Bearer fake-value"))
		})

		Context("invalid config", func() {
			BeforeEach(func() {
				env.Auth.Options = "invalid"
			})

			It("errors", func() {
				runtime.GetEnvironmentReturns(env, nil)

				err := subject.AuthRequest(&req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("parsing authentication options"))
			})
		})
	})
})
