package httpclient_test

import (
	"errors"

	. "github.com/dpb587/ssoca/service/ssh/httpclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dpb587/ssoca/httpclient/httpclientfakes"
	"github.com/dpb587/ssoca/service/ssh/api"
)

var _ = Describe("New", func() {
	It("requires client", func() {
		_, err := New(nil, "fake-service")
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

		subject, err = New(fakeapiclient, "fake-service")
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("GetCAPublicKey", func() {
		Context("request fails", func() {
			It("errors", func() {
				fakeapiclient.APIGetReturns(errors.New("fake-err"))

				_, err := subject.GetCAPublicKey()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Getting"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("works", func() {
			fakeapiclient.APIGetStub = func(_ string, out interface{}) error {
				assertout, ok := out.(*api.CAPublicKeyResponse)
				Expect(ok).To(BeTrue())

				assertout.OpenSSH = "fake-openssh-data"

				return nil
			}

			result, err := subject.GetCAPublicKey()

			Expect(err).ToNot(HaveOccurred())
			Expect(result.OpenSSH).To(Equal("fake-openssh-data"))

			Expect(fakeapiclient.APIGetCallCount()).To(Equal(1))

			path0, _ := fakeapiclient.APIGetArgsForCall(0)
			Expect(path0).To(Equal("/fake-service/ca-public-key"))
		})
	})

	Describe("PostSignPublicKey", func() {
		Context("request fails", func() {
			It("errors", func() {
				fakeapiclient.APIPostReturns(errors.New("fake-err"))

				_, err := subject.PostSignPublicKey(api.SignPublicKeyRequest{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Posting"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("works", func() {
			fakeapiclient.APIPostStub = func(_ string, out interface{}, _ interface{}) error {
				assertout, ok := out.(*api.SignPublicKeyResponse)
				Expect(ok).To(BeTrue())

				assertout.Certificate = "fake-certificate-data"

				return nil
			}

			in := api.SignPublicKeyRequest{
				PublicKey: "fake-public-key-data",
			}

			result, err := subject.PostSignPublicKey(in)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Certificate).To(Equal("fake-certificate-data"))

			Expect(fakeapiclient.APIPostCallCount()).To(Equal(1))

			path0, _, in0 := fakeapiclient.APIPostArgsForCall(0)
			Expect(path0).To(Equal("/fake-service/sign-public-key"))
			Expect(in0).To(Equal(in))
		})
	})
})
