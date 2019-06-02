package username_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
	. "github.com/dpb587/ssoca/auth/authz/filter/username"

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
		Context("single username", func() {
			BeforeEach(func() {
				subject = Requirement{
					Is: "username1",
				}
			})

			It("satisfies when present", func() {
				username := "username1"
				err := subject.VerifyAuthorization(
					&request,
					&auth.Token{
						Attributes: map[auth.TokenAttribute]*string{
							auth.TokenUsernameAttribute: &username,
						},
					},
				)

				Expect(err).ToNot(HaveOccurred())
			})

			It("does not satisfy when not present", func() {
				err := subject.VerifyAuthorization(
					&request,
					&auth.Token{
						ID: "username2",
					},
				)

				Expect(err).To(HaveOccurred())

				aerr, ok := err.(authz.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("username does not match"))
			})

			It("does not satisfy when missing token", func() {
				err := subject.VerifyAuthorization(
					&request,
					nil,
				)

				Expect(err).To(HaveOccurred())

				aerr, ok := err.(authn.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("authentication token missing"))
			})
		})
	})
})
