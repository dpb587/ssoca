package scope_test

import (
	"net/http"

	. "github.com/dpb587/ssoca/authz/filter/scope"

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
		Context("single scope", func() {
			BeforeEach(func() {
				subject = Requirement{
					Present: "scope1",
				}
			})

			It("satisfies when present", func() {
				satisfied, err := subject.IsSatisfied(
					&request,
					&auth.Token{
						Groups: []string{
							"scope1",
						},
					},
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeTrue())
			})

			It("satisfies when more are present", func() {
				satisfied, err := subject.IsSatisfied(
					&request,
					&auth.Token{
						Groups: []string{
							"scope1",
							"scope2",
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
						Groups: []string{
							"scope2",
						},
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
