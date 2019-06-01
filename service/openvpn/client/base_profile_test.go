package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dpb587/ssoca/client/clientfakes"
	"github.com/dpb587/ssoca/httpclient/httpclientfakes"
	. "github.com/dpb587/ssoca/service/openvpn/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseProfile", func() {
	var fakeruntime *clientfakes.FakeRuntime
	var fakeclient *httpclientfakes.FakeClient
	var subject Service

	BeforeEach(func() {
		fakeclient = &httpclientfakes.FakeClient{}

		fakeruntime = &clientfakes.FakeRuntime{}
		fakeruntime.GetAuthInterceptClientReturns(fakeclient, nil)
		fakeruntime.GetClientReturns(fakeclient, nil)

		subject = NewService("fake-name", fakeruntime, nil, nil, nil, nil, nil)
	})

	Describe("unavailable client", func() {
		It("errors", func() {
			fakeruntime.GetAuthInterceptClientReturns(nil, errors.New("fake-err1"))
			fakeruntime.GetClientReturns(nil, errors.New("fake-err1"))

			_, err := subject.BaseProfile(BaseProfileOptions{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Getting client"))
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})
	})

	Describe("failing remote server", func() {
		It("errors", func() {
			fakeclient.GetReturns(nil, errors.New("fake-err1"))

			_, err := subject.BaseProfile(BaseProfileOptions{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})
	})

	It("works", func() {
		fakeclient.GetReturns(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("line1\nline2\n")),
		}, nil)

		profile, err := subject.BaseProfile(BaseProfileOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(profile).To(Equal("line1\nline2\n"))
	})
})
