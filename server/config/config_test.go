package config_test

import (
	. "github.com/dpb587/ssoca/server/config"
	yaml "gopkg.in/yaml.v2"

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

	Describe("ServerTrustedProxies", func() {
		It("parses IPs", func() {
			err := yaml.Unmarshal([]byte(`--- { server: { trusted_proxies: [ 127.0.0.1 ] } }`), &subject)
			Expect(err).ToNot(HaveOccurred())

			Expect(subject.Server.TrustedProxies).To(HaveLen(1))
			Expect(subject.Server.TrustedProxies[0].IP.String()).To(Equal("127.0.0.1"))
			Expect(subject.Server.TrustedProxies[0].Mask.String()).To(Equal("ffffffff"))
		})

		It("parses CIDRs", func() {
			err := yaml.Unmarshal([]byte(`--- { server: { trusted_proxies: [ 127.0.0.0/8 ] } }`), &subject)
			Expect(err).ToNot(HaveOccurred())

			Expect(subject.Server.TrustedProxies).To(HaveLen(1))
			Expect(subject.Server.TrustedProxies[0].IP.String()).To(Equal("127.0.0.0"))
			Expect(subject.Server.TrustedProxies[0].Mask.String()).To(Equal("ff000000"))
		})

		It("converts to IPNet", func() {
			err := yaml.Unmarshal([]byte(`--- { server: { trusted_proxies: [ 127.0.0.0/8 ] } }`), &subject)
			Expect(err).ToNot(HaveOccurred())

			converted := subject.Server.TrustedProxies.AsIPNet()
			Expect(converted).To(HaveLen(1))
			Expect(converted[0].IP.String()).To(Equal("127.0.0.0"))
		})
	})
})
