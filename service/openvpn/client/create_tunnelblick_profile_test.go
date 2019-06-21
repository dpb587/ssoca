package client_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cloudfoundry/bosh-utils/system/fakes"
	"github.com/dpb587/ssoca/client/clientfakes"
	"github.com/dpb587/ssoca/client/config/configfakes"
	"github.com/dpb587/ssoca/httpclient/httpclientfakes"
	. "github.com/dpb587/ssoca/service/openvpn/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateTunnelblickProfile", func() {
	var fakeruntime *clientfakes.FakeRuntime
	var fakefs *fakes.FakeFileSystem
	var fakeclient *httpclientfakes.FakeClient
	var fakeconfigmanager *configfakes.FakeManager
	var subject *Service

	BeforeEach(func() {
		fakefs = fakes.NewFakeFileSystem()
		fakeclient = &httpclientfakes.FakeClient{}

		fakeconfigmanager = &configfakes.FakeManager{}
		fakeconfigmanager.GetSourceReturns("/fake-config/source")

		fakeruntime = &clientfakes.FakeRuntime{}
		fakeruntime.GetAuthInterceptClientReturns(fakeclient, nil)
		fakeruntime.GetClientReturns(fakeclient, nil)
		fakeruntime.GetEnvironmentNameReturns("fake-env")
		fakeruntime.GetConfigManagerReturns(fakeconfigmanager, nil)

		subject = NewService("fake-name", fakeruntime, nil, fakefs, nil, nil)
	})

	Context("failing config manager", func() {
		It("errors", func() {
			fakeruntime.GetConfigManagerReturns(nil, errors.New("fake-err1"))

			_, err := subject.CreateTunnelblickProfile(CreateTunnelblickProfileOpts{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("getting config manager"))
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})
	})

	Describe("resolvable ssoca paths", func() {
		It("errors on invalid paths", func() {
			_, err := subject.CreateTunnelblickProfile(CreateTunnelblickProfileOpts{
				SsocaExec: "nonexistant12345",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("resolving ssoca executable"))
		})
	})

	Describe("failing remote server", func() {
		It("errors", func() {
			fakeclient.GetReturns(nil, errors.New("fake-err1"))

			_, err := subject.CreateTunnelblickProfile(CreateTunnelblickProfileOpts{
				SsocaExec: "false",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("getting base profile"))
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})
	})

	It("works", func() {
		fakeclient.GetReturns(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(strings.NewReader("line1\nline2\n")),
		}, nil)

		profile, err := subject.CreateTunnelblickProfile(CreateTunnelblickProfileOpts{
			Directory: "/fake-tmp1",
			FileName:  "custom1",
			SsocaExec: "env",
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(profile).To(Equal("/fake-tmp1/custom1.tblk"))

		configStat, err := fakefs.Stat("/fake-tmp1/custom1.tblk/config.ovpn")
		Expect(err).ToNot(HaveOccurred())
		Expect(configStat.Mode()).To(Equal(os.FileMode(0400)))

		config, err := fakefs.ReadFileString("/fake-tmp1/custom1.tblk/config.ovpn")
		Expect(err).ToNot(HaveOccurred())
		Expect(config).To(Equal("line1\nline2\n"))

		preConnectStat, err := fakefs.Stat("/fake-tmp1/custom1.tblk/pre-connect.sh")
		Expect(err).ToNot(HaveOccurred())
		Expect(preConnectStat.Mode()).To(Equal(os.FileMode(0500)))

		preConnect, err := fakefs.ReadFileString("/fake-tmp1/custom1.tblk/pre-connect.sh")
		Expect(err).ToNot(HaveOccurred())

		By("passing ssoca executable", func() {
			Expect(preConnect).To(ContainSubstring(`/usr/bin/env --config `))
		})

		By("passing config path", func() {
			Expect(preConnect).To(ContainSubstring(`--config "/fake-config/source"`))
		})

		By("passing environment name", func() {
			Expect(preConnect).To(ContainSubstring(`--environment "fake-env"`))
		})

		By("passing service name", func() {
			Expect(preConnect).To(ContainSubstring(`--service "fake-name"`))
		})
	})
})
