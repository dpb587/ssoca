package httpclient_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/dpb587/ssoca/service/openvpn/httpclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dpb587/ssoca/httpclient/httpclientfakes"
	"github.com/dpb587/ssoca/service/openvpn/api"
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

	Describe("BaseProfile", func() {
		Context("request fails", func() {
			It("errors", func() {
				fakeapiclient.GetReturns(nil, errors.New("fake-err"))

				_, err := subject.BaseProfile()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("getting"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("works", func() {
			fakeapiclient.GetReturns(&http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("line1\nline2\n")),
			}, nil)

			result, err := subject.BaseProfile()

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(`line1
line2
`))

			Expect(fakeapiclient.GetCallCount()).To(Equal(1))

			path0 := fakeapiclient.GetArgsForCall(0)
			Expect(path0).To(Equal("/fake-service/base-profile"))
		})
	})

	Describe("SignUserCSR", func() {
		Context("request fails", func() {
			It("errors", func() {
				fakeapiclient.APIPostReturns(errors.New("fake-err"))

				_, err := subject.SignUserCSR(api.SignUserCSRRequest{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("posting"))
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})

		It("works", func() {
			fakeapiclient.APIPostStub = func(_ string, out interface{}, _ interface{}) error {
				assertout, ok := out.(*api.SignUserCSRResponse)
				Expect(ok).To(BeTrue())

				assertout.Certificate = "fake-certificate-data"
				assertout.Profile = "fake-profile-data"

				return nil
			}

			in := api.SignUserCSRRequest{
				CSR: "fake-csr-data",
			}

			result, err := subject.SignUserCSR(in)

			Expect(err).ToNot(HaveOccurred())
			Expect(result.Certificate).To(Equal("fake-certificate-data"))
			Expect(result.Profile).To(Equal("fake-profile-data"))

			Expect(fakeapiclient.APIPostCallCount()).To(Equal(1))

			path0, _, in0 := fakeapiclient.APIPostArgsForCall(0)
			Expect(path0).To(Equal("/fake-service/sign-user-csr"))
			Expect(in0).To(Equal(in))
		})
	})
})
