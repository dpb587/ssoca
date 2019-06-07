package config_test

import (
	. "github.com/dpb587/ssoca/service/githubauth/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var subject Config

	Describe("ApplyDefaults", func() {
		BeforeEach(func() {
			subject = Config{}
			subject.ApplyDefaults()
		})

		It("defaults auth_url: from https://developer.github.com/v3/oauth/#web-application-flow", func() {
			Expect(subject.AuthURL).To(Equal("https://github.com/login/oauth/authorize"))
		})

		It("defaults token_url: from https://developer.github.com/v3/oauth/#web-application-flow", func() {
			Expect(subject.TokenURL).To(Equal("https://github.com/login/oauth/access_token"))
		})
	})
})
