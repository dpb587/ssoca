package dynamicvalue_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/server/service/dynamicvalue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringValue", func() {
	Describe("Evaluate", func() {
		It("returns", func() {
			value := NewStringValue("something static")

			res, err := value.Evaluate(&http.Request{}, &auth.Token{})

			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("something static"))
		})
	})
})
