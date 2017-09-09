package req_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz/filter/filterfakes"
	. "github.com/dpb587/ssoca/server/service/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RouteHandlerFunc", func() {
	var subject RouteHandlerFunc

	Describe("Route", func() {
		It("returns path", func() {
			subject = RouteHandlerFunc{Path: "/fake-somewhere"}

			Expect(subject.Route()).To(Equal("/fake-somewhere"))
		})
	})

	Describe("Execute", func() {
		It("forwards to handler", func() {
			var called bool
			req := Request{
				RawRequest:  &http.Request{},
				RawResponse: httptest.NewRecorder(),
			}

			subject = RouteHandlerFunc{Func: func(w http.ResponseWriter, r *http.Request) {
				called = true

				Expect(w).To(Equal(req.RawResponse))
				Expect(r).To(Equal(req.RawRequest))
			}}

			subject.Execute(req)

			Expect(called).To(BeTrue())
		})
	})

	Describe("VerifyAuthorization", func() {
		var req *http.Request
		var token *auth.Token

		BeforeEach(func() {
			req = &http.Request{Method: "fake-method"}
			token = &auth.Token{ID: "fake-id"}
		})

		Context("without requirement", func() {
			It("is authorized", func() {
				err := subject.VerifyAuthorization(req, token)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with requirement", func() {
			var requirement *filterfakes.FakeRequirement

			BeforeEach(func() {
				requirement = &filterfakes.FakeRequirement{}
				requirement.VerifyAuthorizationReturns(errors.New("fake-err1"))

				subject.Requirement = requirement
			})

			It("delegates", func() {
				err := subject.VerifyAuthorization(req, token)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err1"))

				Expect(requirement.VerifyAuthorizationCallCount()).To(Equal(1))
				arg1, arg2 := requirement.VerifyAuthorizationArgsForCall(0)
				Expect(arg1).To(Equal(req))
				Expect(arg2).To(Equal(token))
			})
		})
	})
})
