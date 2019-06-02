package public_test

import (
	"net/http"

	. "github.com/dpb587/ssoca/auth/authz/filter/public"

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

		It("satisfies without token", func() {
			err := subject.VerifyAuthorization(
				&request,
				nil,
			)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
