package authenticated_test

import (
	"net/http"

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

	Describe("IsSatisfied", func() {
		BeforeEach(func() {
			subject = Requirement{}
		})

		It("satisfies with token", func() {
			satisfied, err := subject.IsSatisfied(
				&request,
				&auth.Token{ID: "test"},
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(satisfied).To(BeTrue())
		})

		It("does not satisfy without token", func() {
			satisfied, err := subject.IsSatisfied(
				&request,
				nil,
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(satisfied).To(BeFalse())
		})
	})
})
