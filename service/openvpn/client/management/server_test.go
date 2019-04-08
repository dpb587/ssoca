package management_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/dpb587/ssoca/internal/nettest"
	. "github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/management/managementfakes"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var subject Client
	var fakeHandler *managementfakes.FakeServerHandler
	var loggerOutput *bytes.Buffer
	var localConn, remoteConn net.Conn

	BeforeEach(func() {
		fakeHandler = &managementfakes.FakeServerHandler{}
		localConn, remoteConn = nettest.NewConnection()

		loggerOutput = &bytes.Buffer{}

		logger := logrus.New()
		logger.Level = logrus.DebugLevel
		logger.Out = loggerOutput
		logger.Formatter = &logrus.JSONFormatter{}

		subject = NewClient(remoteConn, fakeHandler, "", logger)

		go subject.Run()
	})

	AfterEach(func() {
		localConn.Close()
		remoteConn.Close()
	})

	Describe("realtime messages", func() {
		Context("INFO", func() {
			It("emits a log message", func() {
				_, err := localConn.Write([]byte(">INFO:fake-info-message\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() string { return loggerOutput.String() }).Should(ContainSubstring(`"level":"info","msg":"fake-info-message"`))
			})
		})

		Context("NEED-CERTIFICATE", func() {
			It("would provide a certificate", func() {
				callback := make(chan struct{}, 0)

				fakeHandler.NeedCertificateReturns(
					func(_ io.Writer, data string) (ServerHandlerCallback, error) {
						defer GinkgoRecover()

						Expect(data).To(Equal("SUCCESS\n"))
						close(callback)

						return nil, nil
					},
					nil,
				)

				_, err := localConn.Write([]byte(">NEED-CERTIFICATE:fake-id\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return fakeHandler.NeedCertificateCallCount() }).Should(Equal(1))

				arg1, arg2 := fakeHandler.NeedCertificateArgsForCall(0)
				Expect(arg1).To(Equal(remoteConn))
				Expect(arg2).To(Equal("fake-id"))

				_, err = localConn.Write([]byte("SUCCESS\n"))
				Expect(err).ToNot(HaveOccurred())

				select {
				case <-callback:
				case <-time.Tick(time.Second):
					Fail("Expected callback to be executed")
				}
			})

			It("causes SIGTERM when renwewal fails", func() {
				fakeHandler.NeedCertificateReturns(nil, errors.New("fake-err1"))

				_, err := localConn.Write([]byte(">NEED-CERTIFICATE:fake-id\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return fakeHandler.NeedCertificateCallCount() }).Should(Equal(1))

				remoteConn.Close()

				all, err := ioutil.ReadAll(localConn)
				Expect(err).ToNot(HaveOccurred())

				Expect(all).To(Equal([]byte("signal SIGTERM\n")))
			})
		})

		Context("RSA_SIGN", func() {
			It("would sign a certificate", func() {
				callback := make(chan struct{}, 0)

				fakeHandler.SignRSAReturns(
					func(_ io.Writer, data string) (ServerHandlerCallback, error) {
						defer GinkgoRecover()

						Expect(data).To(Equal("SUCCESS\n"))
						close(callback)

						return nil, nil
					},
					nil,
				)

				_, err := localConn.Write([]byte(">RSA_SIGN:fake-data\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() int { return fakeHandler.SignRSACallCount() }).Should(Equal(1))

				arg1, arg2 := fakeHandler.SignRSAArgsForCall(0)
				Expect(arg1).To(Equal(remoteConn))
				Expect(arg2).To(Equal("fake-data"))

				_, err = localConn.Write([]byte("SUCCESS\n"))
				Expect(err).ToNot(HaveOccurred())

				select {
				case <-callback:
				case <-time.Tick(time.Second):
					Fail("Expected callback to be executed")
				}
			})
		})

		Context("FATAL", func() {
			It("would exit", func() {
				_, err := localConn.Write([]byte(">FATAL:fake-fatal-message\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() string { return loggerOutput.String() }).Should(ContainSubstring(`"level":"error","msg":"fake-fatal-message"`))
			})
		})

		Context("unrecognized", func() {
			It("warns", func() {
				_, err := localConn.Write([]byte(">MISSING-CMD:fake-data\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() string { return loggerOutput.String() }).Should(ContainSubstring(`"level":"warning","msg":"unexpected realtime message: MISSING-CMD"`))
			})
		})
	})

	Describe("callback messages", func() {
		Context("unrecognized", func() {
			It("warns", func() {
				_, err := localConn.Write([]byte("UNKNOWN\n"))
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() string { return loggerOutput.String() }).Should(ContainSubstring(`"level":"warning","msg":"unexpected callback message: UNKNOWN\n"`))
			})
		})
	})
})
