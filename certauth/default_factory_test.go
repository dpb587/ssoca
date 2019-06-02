package certauth_test

import (
	"errors"

	. "github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultFactory", func() {
	Describe("Register", func() {
		It("works", func() {
			provider := &certauthfakes.FakeProvider{}
			providerOptions := map[string]interface{}{"key1": "val1"}

			providerFactory := &certauthfakes.FakeProviderFactory{}
			providerFactory.CreateReturns(provider, nil)

			factory := NewDefaultFactory()
			factory.Register("fake1", providerFactory)

			prv, err := factory.Create("name1", "fake1", providerOptions)

			Expect(err).ToNot(HaveOccurred())
			Expect(prv).To(Equal(provider))

			arg0, arg1 := providerFactory.CreateArgsForCall(0)
			Expect(arg0).To(Equal("name1"))
			Expect(arg1).To(Equal(providerOptions))
		})

		Context("unknown provider", func() {
			It("errors", func() {
				factory := NewDefaultFactory()

				_, err := factory.Create("name1", "unknown1", map[string]interface{}{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown1"))
			})
		})

		Context("factory errors", func() {
			It("errors", func() {
				providerFactory := &certauthfakes.FakeProviderFactory{}
				providerFactory.CreateReturns(nil, errors.New("fake-error1"))

				factory := NewDefaultFactory()
				factory.Register("fake1", providerFactory)

				_, err := factory.Create("name1", "fake1", map[string]interface{}{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("creating provider fake1"))
			})
		})
	})
})
