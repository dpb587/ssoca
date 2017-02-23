package config_test

import (
	. "github.com/dpb587/ssoca/server/config"

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

		It("defaults server host to all interfaces", func() {
			Expect(subject.Server.Host).To(Equal("0.0.0.0"))
		})

		It("defaults server port", func() {
			Expect(subject.Server.Port).To(Equal(18705))
		})
	})
})
