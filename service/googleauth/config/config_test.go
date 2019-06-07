package config_test

import (
	. "github.com/dpb587/ssoca/service/googleauth/config"

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

		It("defaults auth_url", func() {
			Expect(subject.AuthURL).To(Equal("https://accounts.google.com/o/oauth2/v2/auth"))
		})

		It("defaults token_url", func() {
			Expect(subject.TokenURL).To(Equal("https://www.googleapis.com/oauth2/v4/token"))
		})
	})
})
