package config_test

import (
	. "github.com/dpb587/ssoca/service/githubauth/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var subject Config

	BeforeEach(func() {
		subject = Config{}
	})

	Describe("ApplyDefaults", func() {
		BeforeEach(func() {
			subject.ApplyDefaults()
		})

		It("defaults auth_url: from https://developer.github.com/v3/oauth/#web-application-flow", func() {
			Expect(subject.AuthURL).To(Equal("https://github.com/login/oauth/authorize"))
		})

		It("defaults token_url: from https://developer.github.com/v3/oauth/#web-application-flow", func() {
			Expect(subject.TokenURL).To(Equal("https://github.com/login/oauth/access_token"))
		})
	})

	Describe("ApplyRedirectDefaults", func() {
		BeforeEach(func() {
			subject.ApplyDefaults()
		})

		It("active", func() {
			subject.ApplyRedirectDefaults("http://success", "http://failure")
			Expect(subject.FailureRedirectURL).To(Equal("http://failure"))
			Expect(subject.SuccessRedirectURL).To(Equal("http://success"))
		})

		It("inactive", func() {
			subject.FailureRedirectURL = "http://existing-failure"
			subject.SuccessRedirectURL = "http://existing-success"

			subject.ApplyRedirectDefaults("http://success", "http://failure")
			Expect(subject.FailureRedirectURL).To(Equal("http://existing-failure"))
			Expect(subject.SuccessRedirectURL).To(Equal("http://existing-success"))
		})
	})
})
