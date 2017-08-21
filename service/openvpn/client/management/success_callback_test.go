package management_test

import (
	. "github.com/dpb587/ssoca/service/openvpn/client/management"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuccessCallback", func() {
	It("successful response", func() {
		rcb, err := SuccessCallback(nil, "SUCCESS: data")
		Expect(err).ToNot(HaveOccurred())
		Expect(rcb).To(BeNil())
	})

	It("failure response", func() {
		_, err := SuccessCallback(nil, "FAILURE: data")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Bad management command result"))
	})
})
