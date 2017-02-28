package api_test

import (
	"errors"
	"io/ioutil"
	"strings"

	. "github.com/dpb587/ssoca/server/api"

	"net/http"
	"net/http/httptest"

	"github.com/Sirupsen/logrus"
	logrustest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/internal/internalfakes"
	"github.com/dpb587/ssoca/server/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var res httptest.ResponseRecorder
	var req *http.Request
	var authService *servicefakes.FakeAuthService
	var apiService *servicefakes.FakeService
	var logger logrus.FieldLogger

	BeforeEach(func() {
		res = *httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://localhost/file?name=test1", nil)

		authService = &servicefakes.FakeAuthService{}
		apiService = &servicefakes.FakeService{}
		apiService.IsAuthorizedReturns(true, nil)

		logger, _ = logrustest.NewNullLogger()
	})

	Describe("CreateHandler", func() {
		Context("handler inputs", func() {
			Context("http.ResponseWriter", func() {
				It("works", func() {
					handler := &internalfakes.FakeFakeInHttpResponseWriter{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(handler.ExecuteArgsForCall(0)).To(Equal(&res))
				})
			})

			Context("*http.Request", func() {
				It("works", func() {
					handler := &internalfakes.FakeFakeInHttpRequest{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(handler.ExecuteArgsForCall(0).URL).To(Equal(req.URL))
				})
			})

			Context("auth.Token", func() {
				It("works", func() {
					var authToken = auth.Token{ID: "fake-id"}

					authService.ParseRequestAuthReturns(&authToken, nil)

					handler := &internalfakes.FakeFakeInAuthToken{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(handler.ExecuteArgsForCall(0)).To(Equal(&authToken))
				})

				Context("missing token", func() {
					It("errors with 401", func() {
						handler := &internalfakes.FakeFakeInAuthToken{}

						wrapper, err := CreateHandler(authService, apiService, handler, logger)

						Expect(err).ToNot(HaveOccurred())

						wrapper.ServeHTTP(&res, req)

						Expect(handler.ExecuteCallCount()).To(Equal(0))
						Expect(res.Code).To(Equal(401))
					})
				})
			})

			Context("api payload", func() {
				var handled bool

				BeforeEach(func() {
					handled = false
				})

				It("works", func() {
					handler := &internalfakes.FakeFakeInApiPayload{}

					req.Body = ioutil.NopCloser(strings.NewReader(`{"test":true}`))

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(handler.ExecuteArgsForCall(0)).To(BeEquivalentTo(map[string]interface{}{"test": true}))
				})

				Context("invalid json", func() {
					It("errors", func() {
						handler := &internalfakes.FakeFakeInApiPayload{}

						req.Body = ioutil.NopCloser(strings.NewReader(`{"test"`))

						wrapper, err := CreateHandler(authService, apiService, handler, logger)

						Expect(err).ToNot(HaveOccurred())

						wrapper.ServeHTTP(&res, req)

						Expect(handler.ExecuteCallCount()).To(Equal(0))
						Expect(res.Code).To(Equal(400))
					})
				})
			})
		})

		Context("handler outputs", func() {
			Context("return", func() {
				It("works", func() {
					handler := &internalfakes.FakeFake{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
				})
			})

			Context("return error", func() {
				It("works with error", func() {
					handler := &internalfakes.FakeFakeOutError{}
					handler.ExecuteReturns(errors.New("fake-err"))

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(res.Code).To(Equal(500))
				})

				It("works without error", func() {
					handler := &internalfakes.FakeFakeOutError{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(res.Code).ToNot(Equal(500))
				})
			})

			Context("return interface{}, error", func() {
				var handled bool

				BeforeEach(func() {
					handled = false
				})

				It("works with json", func() {
					handler := &internalfakes.FakeFakeOutInterfaceError{}
					handler.ExecuteReturns(map[string]interface{}{"output": true}, nil)

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(res.Body.String()).To(Equal(`{
  "output": true
}
`))
				})

				It("works with errors", func() {
					handler := &internalfakes.FakeFakeOutInterfaceError{}
					handler.ExecuteReturns(nil, errors.New("fake-err"))

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(1))
					Expect(res.Code).To(Equal(500))
				})
			})

			Context("return something else", func() {
				It("errors", func() {
					handler := &internalfakes.FakeFakeOutOther{}

					_, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid handler function"))
				})
			})
		})
	})

	Describe("ServeHTTP", func() {
		Context("when auth service fails", func() {
			It("fails the request", func() {
				authService.ParseRequestAuthReturns(nil, errors.New("fake-err"))

				handler := &internalfakes.FakeFake{}

				wrapper, err := CreateHandler(authService, apiService, handler, logger)

				Expect(err).ToNot(HaveOccurred())

				wrapper.ServeHTTP(&res, req)

				Expect(handler.ExecuteCallCount()).To(Equal(0))
				Expect(res.Code).To(Equal(500))
			})
		})

		Context("when authorization fails", func() {
			It("fails the request", func() {
				apiService.IsAuthorizedReturns(false, errors.New("fake-err"))

				handler := &internalfakes.FakeFake{}

				wrapper, err := CreateHandler(authService, apiService, handler, logger)

				Expect(err).ToNot(HaveOccurred())

				wrapper.ServeHTTP(&res, req)

				Expect(handler.ExecuteCallCount()).To(Equal(0))
				Expect(res.Code).To(Equal(500))
			})
		})

		Context("when authorization rejected", func() {
			It("errors with 403", func() {
				apiService.IsAuthorizedReturns(false, nil)

				handler := &internalfakes.FakeFake{}

				wrapper, err := CreateHandler(authService, apiService, handler, logger)

				Expect(err).ToNot(HaveOccurred())

				wrapper.ServeHTTP(&res, req)

				Expect(handler.ExecuteCallCount()).To(Equal(0))
				Expect(res.Code).To(Equal(403))
			})
		})
	})
})
