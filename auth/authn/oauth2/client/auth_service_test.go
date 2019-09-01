package client_test

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/dpb587/ssoca/auth/authn"
	. "github.com/dpb587/ssoca/auth/authn/oauth2/client"

	"github.com/dpb587/ssoca/client/clientfakes"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/config/configfakes"
	"github.com/pkg/errors"

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
	var stdout, stderr *bytes.Buffer

	var runtime *clientfakes.FakeRuntime
	var configManager *configfakes.FakeManager

	BeforeEach(func() {
		env = config.EnvironmentState{
			Auth: &config.EnvironmentAuthState{},
			URL:  "https://192.0.2.99",
		}

		stderr = bytes.NewBuffer(nil)
		stdout = bytes.NewBuffer(nil)

		ui = &uifakes.FakeUI{}
		cmdRunner = &boshsysfakes.FakeCmdRunner{}

		configManager = &configfakes.FakeManager{}

		runtime = &clientfakes.FakeRuntime{}
		runtime.GetUIReturns(ui)
		runtime.GetStdinReturns(bytes.NewBuffer(nil))
		runtime.GetStderrReturns(stderr)
		runtime.GetStdoutReturns(stdout)
		runtime.GetConfigManagerReturns(configManager, nil)

		subject = NewAuthService("fake-name", "fake-type", runtime, cmdRunner)
	})

	Describe("AuthLogin", func() {
		It("shows login url and retrieves token from user", func() {
			runtime.GetEnvironmentReturns(env, nil)
			runtime.GetStdinReturns(bytes.NewBufferString("mytoken\r\n"))

			err := subject.AuthLogin(context.Background())

			Expect(err).ToNot(HaveOccurred())

			Expect(configManager.SetEnvironmentCallCount()).To(Equal(1))

			envState := configManager.SetEnvironmentArgsForCall(0)
			token := envState.Auth.Options

			authConfig, ok := token.(authn.AuthorizationToken)
			Expect(ok).To(BeTrue())

			Expect(authConfig.Type).To(Equal("Bearer"))
			Expect(authConfig.Value).To(Equal("mytoken"))

			Expect(stderr.String()).To(ContainSubstring("https://192.0.2.99/fake-name/initiate"))
		})

		It("respects timeouts", func() {
			runtime.GetEnvironmentReturns(env, nil)

			ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second)
			defer ctxCancel()

			err := subject.AuthLogin(ctx)

			Expect(err).To(HaveOccurred())
			Expect(errors.Cause(err)).To(Equal(context.DeadlineExceeded))
		})
	})

	Describe("AuthLogout", func() {
		It("does not error", func() {
			err := subject.AuthLogout(context.Background())

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
