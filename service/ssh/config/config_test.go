package config_test

import (
	. "github.com/dpb587/ssoca/service/ssh/config"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var subject Config

	BeforeEach(func() {
		subject = Config{}
	})

	// Describe("ApplyDefaults", func() {
	// 	Describe("with defaults", func() {
	// 		BeforeEach(func() {
	// 			subject.ApplyDefaults()
	// 		})
	//
	// 		It("defaults extensions", func() {
	// 			Expect(subject.Extensions).To(Equal(ExtensionDefaults))
	// 		})
	// 	})
	//
	// 	Describe("not really dealing with defaults but its convenient so...", func() {
	// 		It("ignores extensions if ssoca-no-defaults is used", func() {
	// 			subject.Extensions = Extensions{ExtensionNoDefaults}
	// 			subject.ApplyDefaults()
	//
	// 			Expect(subject.Extensions).To(HaveLen(0))
	// 		})
	// 	})
	// })
})
