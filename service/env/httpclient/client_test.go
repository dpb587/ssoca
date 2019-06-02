package httpclient_test

import (
	"errors"

	. "github.com/dpb587/ssoca/service/env/httpclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dpb587/ssoca/httpclient/httpclientfakes"
	"github.com/dpb587/ssoca/service/env/api"
)

var _ = Describe("New", func() {
	It("requires client", func() {
		_, err := New(nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("client is nil"))
	})
})

var _ = Describe("Client", func() {
	var subject Client
	var fakeapiclient *httpclientfakes.FakeClient

	BeforeEach(func() {
		var err error

		fakeapiclient = &httpclientfakes.FakeClient{}

		subject, err = New(fakeapiclient)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("GetInfo", func() {
		Context("request fails", func() {
			It("errors", func() {
				fakeapiclient.APIGetReturns(errors.New("fake-err"))

				_, err := subject.GetInfo()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("getting"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("works", func() {
			fakeapiclient.APIGetStub = func(_ string, out interface{}) error {
				assertout, ok := out.(*api.InfoResponse)
				Expect(ok).To(BeTrue())

				assertout.Version = "fake-version-data"

				return nil
			}

			result, err := subject.GetInfo()

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Version).To(Equal("fake-version-data"))

			Expect(fakeapiclient.APIGetCallCount()).To(Equal(1))

			path0, _ := fakeapiclient.APIGetArgsForCall(0)
			Expect(path0).To(Equal("/env/info"))
		})
	})
})
