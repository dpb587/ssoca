package config_test

import (
	. "github.com/dpb587/ssoca/service/openvpn/config"

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

		It("defaults certauth: default", func() {
			Expect(subject.CertAuthName).To(Equal("default"))
		})

		It("defaults validity", func() {
			Expect(subject.ValidityString).To(Equal("2m"))
		})
	})
})
