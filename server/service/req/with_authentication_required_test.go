package req_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authn"
	. "github.com/dpb587/ssoca/server/service/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithAuthenticationRequired", func() {
	var subject WithAuthenticationRequired

	Describe("VerifyAuthorization", func() {
		Context("with authorization token", func() {
			It("is authorized", func() {
				err := subject.VerifyAuthorization(&http.Request{}, &auth.Token{ID: "authenticated"})

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("without authorization token", func() {
			It("is not authorized", func() {
				err := subject.VerifyAuthorization(&http.Request{}, nil)

				Expect(err).To(HaveOccurred())

				aerr, ok := err.(authn.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("authentication token missing"))
			})
		})
	})
})
