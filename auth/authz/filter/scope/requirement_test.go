package scope_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
	. "github.com/dpb587/ssoca/auth/authz/filter/scope"

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
		Context("single scope", func() {
			BeforeEach(func() {
				subject = Requirement{
					Present: "scope1",
				}
			})

			It("satisfies when present", func() {
				err := subject.VerifyAuthorization(
					&request,
					&auth.Token{
						Groups: []string{
							"scope1",
						},
					},
				)

				Expect(err).ToNot(HaveOccurred())
			})

			It("satisfies when more are present", func() {
				err := subject.VerifyAuthorization(
					&request,
					&auth.Token{
						Groups: []string{
							"scope1",
							"scope2",
						},
					},
				)

				Expect(err).ToNot(HaveOccurred())
			})

			It("does not satisfy when not present", func() {
				err := subject.VerifyAuthorization(
					&request,
					&auth.Token{
						Groups: []string{
							"scope2",
						},
					},
				)

				aerr, ok := err.(authz.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("scope is missing"))
			})

			It("does not satisfy when missing token", func() {
				err := subject.VerifyAuthorization(
					&request,
					nil,
				)

				aerr, ok := err.(authn.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("authentication token missing"))
			})
		})
	})
})
