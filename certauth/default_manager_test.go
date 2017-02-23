package certauth_test

import (
	. "github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultManager", func() {
	var subject Manager

	BeforeEach(func() {
		subject = NewDefaultManager()
	})

	Describe("Get", func() {
		It("retrieves providers", func() {
			provider := &certauthfakes.FakeProvider{}
			provider.NameReturns("test1")

			subject.Add(provider)

			get, err := subject.Get("test1")

			Expect(err).ToNot(HaveOccurred())
			Expect(get).To(Equal(provider))
		})

		Context("non-existant provider", func() {
			It("errors", func() {
				_, err := subject.Get("test1")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unrecognized name: test1"))
			})
		})
	})
})
