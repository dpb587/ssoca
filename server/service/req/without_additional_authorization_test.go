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

	Describe("IsAuthorized", func() {
		Context("with authorization token", func() {
			It("is authorized", func() {
				authz, err := subject.IsAuthorized(&http.Request{}, &auth.Token{ID: "authenticated"})

				Expect(err).ToNot(HaveOccurred())
				Expect(authz).To(BeTrue())
			})
		})

		Context("without authorization token", func() {
			It("is authorized", func() {
				authz, err := subject.IsAuthorized(&http.Request{}, nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(authz).To(BeTrue())
			})
		})
	})
})
