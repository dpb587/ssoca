package req_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/server/service/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithoutAdditionalAuthorization", func() {
	var subject WithoutAdditionalAuthorization

	Describe("VerifyAuthorization", func() {
		Context("with authorization token", func() {
			It("is authorized", func() {
				err := subject.VerifyAuthorization(&http.Request{}, &auth.Token{ID: "authenticated"})
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("without authorization token", func() {
			It("is authorized", func() {
				err := subject.VerifyAuthorization(&http.Request{}, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
