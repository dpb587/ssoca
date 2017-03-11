package username_test

import (
	"net/http"

	. "github.com/dpb587/ssoca/authz/filter/username"

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

	Describe("IsSatisfied", func() {
		Context("single username", func() {
			BeforeEach(func() {
				subject = Requirement{
					Is: "username1",
				}
			})

			It("satisfies when present", func() {
				username := "username1"
				satisfied, err := subject.IsSatisfied(
					&request,
					&auth.Token{
						Attributes: map[auth.TokenAttribute]*string{
							auth.TokenUsernameAttribute: &username,
						},
					},
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeTrue())
			})

			It("does not satisfy when not present", func() {
				satisfied, err := subject.IsSatisfied(
					&request,
					&auth.Token{
						ID: "username2",
					},
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeFalse())
			})

			It("does not satisfy when missing token", func() {
				satisfied, err := subject.IsSatisfied(
					&request,
					nil,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeFalse())
			})
		})
	})
})
