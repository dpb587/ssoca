package authenticated_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth/authn"
	. "github.com/dpb587/ssoca/auth/authz/filter/authenticated"

	"github.com/dpb587/ssoca/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requirement", func() {
	var request http.Request
	var subject Requirement

	BeforeEach(func() {
		request = http.Request{}
	})

	Describe("VerifyAuthorization", func() {
		BeforeEach(func() {
			subject = Requirement{}
		})

		It("satisfies with token", func() {
			err := subject.VerifyAuthorization(
				&request,
				&auth.Token{ID: "test"},
			)

			Expect(err).ToNot(HaveOccurred())
		})

		It("does not satisfy without username", func() {
			err := subject.VerifyAuthorization(
				&request,
				&auth.Token{},
			)

			Expect(err).To(HaveOccurred())

			aerr, ok := err.(authn.Error)
			Expect(ok).To(BeTrue())
			Expect(aerr.Error()).To(Equal("Authentication ID missing"))
		})

		It("does not satisfy without token", func() {
			err := subject.VerifyAuthorization(
				&request,
				nil,
			)

			Expect(err).To(HaveOccurred())

			aerr, ok := err.(authn.Error)
			Expect(ok).To(BeTrue())
			Expect(aerr.Error()).To(Equal("Authentication token missing"))
		})
	})
})
