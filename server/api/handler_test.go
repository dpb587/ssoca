package api_test

import (
	"errors"

	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
	. "github.com/dpb587/ssoca/server/api"

	"net/http"
	"net/http/httptest"

	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req/reqfakes"
	"github.com/dpb587/ssoca/server/service/servicefakes"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"

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

		logger, _ = logrustest.NewNullLogger()
	})

	Describe("ServeHTTP", func() {
		Context("when auth service fails", func() {
			It("fails the request", func() {
				authService.ParseRequestAuthReturns(nil, errors.New("fake-err"))

				handler := &reqfakes.FakeRouteHandler{}

				wrapper, err := CreateHandler(authService, apiService, handler, logger)

				Expect(err).ToNot(HaveOccurred())

				wrapper.ServeHTTP(&res, req)

				Expect(handler.ExecuteCallCount()).To(Equal(0))
				Expect(res.Code).To(Equal(500))
			})
		})

		Context("authorization", func() {
			Context("authentication errors", func() {
				It("errors with 401", func() {
					apiService.VerifyAuthorizationReturns(authn.NewError(errors.New("fake-err")))

					handler := &reqfakes.FakeRouteHandler{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(0))
					Expect(res.Code).To(Equal(401))
				})
			})

			Context("authorization errors", func() {
				It("errors with 403", func() {
					apiService.VerifyAuthorizationReturns(authz.NewError(errors.New("fake-err")))

					handler := &reqfakes.FakeRouteHandler{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(0))
					Expect(res.Code).To(Equal(403))
				})
			})

			Context("generic errors", func() {
				It("errors with 500", func() {
					apiService.VerifyAuthorizationReturns(errors.New("fake-err"))

					handler := &reqfakes.FakeRouteHandler{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(0))
					Expect(res.Code).To(Equal(500))
				})
			})

			Context("custom status codes", func() {
				It("errors with the custom status", func() {
					apiService.VerifyAuthorizationReturns(apierr.NewError(errors.New("fake-err"), 418, ""))

					handler := &reqfakes.FakeRouteHandler{}

					wrapper, err := CreateHandler(authService, apiService, handler, logger)

					Expect(err).ToNot(HaveOccurred())

					wrapper.ServeHTTP(&res, req)

					Expect(handler.ExecuteCallCount()).To(Equal(0))
					Expect(res.Code).To(Equal(418))
					Expect(res.Body.String()).To(ContainSubstring("I'm a teapot"))
				})
			})
		})

		Context("when handler fails", func() {
			It("errors with 500", func() {
				handler := &reqfakes.FakeRouteHandler{}
				handler.ExecuteReturns(errors.New("fake-err1"))

				wrapper, err := CreateHandler(authService, apiService, handler, logger)

				Expect(err).ToNot(HaveOccurred())

				wrapper.ServeHTTP(&res, req)

				Expect(handler.ExecuteCallCount()).To(Equal(1))
				Expect(res.Code).To(Equal(500))
				Expect(res.Body.String()).To(Equal(`{
  "error": {
    "message": "Internal Server Error",
    "status": 500
  }
}
`))
			})
		})
	})
})
