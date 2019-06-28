package server_test

import (
	"net/http"

	serverservice "github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/server/service/req/reqfakes"
	servicefakes "github.com/dpb587/ssoca/server/service/servicefakes"
	"github.com/dpb587/ssoca/service"
	. "github.com/dpb587/ssoca/service/auth/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var subject *Service

	Describe("interface", func() {
		It("github.com/dpb587/ssoca/server/service.Service", func() {
			var _ serverservice.Service = (*Service)(nil)
		})
	})

	Context("basics", func() {
		var upstreamauth *servicefakes.FakeAuthService

		BeforeEach(func() {
			upstreamauth = &servicefakes.FakeAuthService{}

			subject = NewService(upstreamauth)
		})

		Describe("Name", func() {
			BeforeEach(func() {
				upstreamauth.NameReturns("fake-upstream-name")
			})

			It("works", func() {
				Expect(subject.Name()).To(Equal("auth"))
			})
		})

		Describe("Type", func() {
			BeforeEach(func() {
				upstreamauth.TypeReturns("fake-upstream-type")
			})

			It("works", func() {
				Expect(subject.Type()).To(Equal(service.Type("fake-upstream-type")))
			})
		})

		Describe("Metadata", func() {
			BeforeEach(func() {
				upstreamauth.MetadataReturns("fake-upstream-metadata")
			})

			It("works", func() {
				Expect(subject.Metadata()).To(Equal("fake-upstream-metadata"))
			})
		})

		Describe("GetRoutes", func() {
			BeforeEach(func() {
				upstreamauth.GetRoutesReturns([]req.RouteHandler{&reqfakes.FakeRouteHandler{}})
			})

			It("delegates", func() {
				routes := subject.GetRoutes()
				Expect(routes).To(HaveLen(2))
				Expect(upstreamauth.GetRoutesCallCount()).To(Equal(1))
			})
		})

		Describe("SupportsRequestAuth", func() {
			It("delegates", func() {
				_, _ = subject.SupportsRequestAuth(http.Request{})
				Expect(upstreamauth.SupportsRequestAuthCallCount()).To(Equal(1))
			})
		})

		Describe("SupportsRequestAuth", func() {
			It("delegates", func() {
				_, _ = subject.SupportsRequestAuth(http.Request{})
				Expect(upstreamauth.SupportsRequestAuthCallCount()).To(Equal(1))
			})
		})

		Describe("ParseRequestAuth", func() {
			It("delegates", func() {
				_, _ = subject.ParseRequestAuth(http.Request{})
				Expect(upstreamauth.ParseRequestAuthCallCount()).To(Equal(1))
			})
		})

		Describe("VerifyAuthorization", func() {
			It("delegates", func() {
				_ = subject.VerifyAuthorization(http.Request{}, nil)
				Expect(upstreamauth.VerifyAuthorizationCallCount()).To(Equal(1))
			})
		})
	})
})
