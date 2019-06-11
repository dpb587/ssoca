package client_test

import (
	"errors"
	"net/http"
	"strings"

	. "github.com/dpb587/ssoca/service/githubauth/client"

	"github.com/dpb587/ssoca/client/clientfakes"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/service/env/api"

	uifakes "github.com/cloudfoundry/bosh-cli/ui/fakes"
	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var subject *Service

	var env config.EnvironmentState
	var ui uifakes.FakeUI
	var cmdRunner boshsysfakes.FakeCmdRunner

	var runtime clientfakes.FakeRuntime

	BeforeEach(func() {
		env = config.EnvironmentState{
			Auth: &config.EnvironmentAuthState{},
			URL:  "https://192.0.2.99",
		}

		ui = uifakes.FakeUI{}
		cmdRunner = boshsysfakes.FakeCmdRunner{}

		runtime = clientfakes.FakeRuntime{}
		runtime.GetUIReturns(&ui)

		subject = NewService("fake-name", &runtime, &cmdRunner)
	})

	Describe("AuthLogin", func() {
		XIt("shows login url and retrieves token from user", func() {
			runtime.GetEnvironmentReturns(env, nil)

			ui.AskedText = append(ui.AskedText, uifakes.Answer{
				Text:  "mytoken",
				Error: nil,
			})

			token, err := subject.AuthLogin(api.InfoServiceResponse{})

			Expect(err).ToNot(HaveOccurred())

			authConfig, ok := token.(AuthConfig)
			Expect(ok).To(BeTrue())

			Expect(authConfig.Token).To(Equal("mytoken"))

			Expect(strings.Join(ui.Said, "")).To(ContainSubstring("https://192.0.2.99/auth/initiate"))
		})

		XIt("propagates errors on input fail", func() {
			runtime.GetEnvironmentReturns(env, nil)

			ui.AskedText = append(ui.AskedText, uifakes.Answer{
				Error: errors.New("ctrl c"),
			})

			_, err := subject.AuthLogin(api.InfoServiceResponse{})

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

		It("adds authorization header", func() {
			env.Auth.Options = AuthConfig{
				Token: "mytoken",
			}

			runtime.GetEnvironmentReturns(env, nil)

			err := subject.AuthRequest(&req)

			Expect(err).ToNot(HaveOccurred())
			Expect(req.Header.Get("Authorization")).To(Equal("bearer mytoken"))
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
