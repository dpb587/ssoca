package authorized_test

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz/filter/filterfakes"
	. "github.com/dpb587/ssoca/auth/authz/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/server/service/req/reqfakes"
	"github.com/dpb587/ssoca/server/service/servicefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var service servicefakes.FakeService
	var requirement filterfakes.FakeRequirement

	var subject Service

	BeforeEach(func() {
		service = servicefakes.FakeService{}
		requirement = filterfakes.FakeRequirement{}

		subject = NewService(&service, &requirement)
	})

	Describe("Type", func() {
		It("delegates to service", func() {
			service.TypeReturns("fake-type")
			Expect(subject.Type()).To(Equal("fake-type"))
		})
	})

	Describe("Version", func() {
		It("delegates to service", func() {
			service.VersionReturns("fake-version")
			Expect(subject.Version()).To(Equal("fake-version"))
		})
	})

	Describe("Name", func() {
		It("delegates to service", func() {
			service.NameReturns("fake-name")
			Expect(subject.Name()).To(Equal("fake-name"))
		})
	})

	Describe("Metadata", func() {
		It("delegates to service", func() {
			metadata := "something"

			service.MetadataReturns(metadata)
			Expect(subject.Metadata()).To(Equal(metadata))
		})
	})

	Describe("GetRoutes", func() {
		It("delegates to service", func() {
			routes := []req.RouteHandler{&reqfakes.FakeRouteHandler{}}

			service.GetRoutesReturns(routes)

			Expect(subject.GetRoutes()).To(Equal(routes))
		})
	})

	Describe("VerifyAuthorization", func() {
		var req http.Request
		var token *auth.Token

		BeforeEach(func() {
			req = http.Request{}
			token = &auth.Token{}
		})

		Context("requirement fails", func() {
			It("is not authorized and does not invoke service authorization check", func() {
				requirement.VerifyAuthorizationReturns(errors.New("fake-err1"))

				err := subject.VerifyAuthorization(req, token)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("fake-err1"))

				Expect(service.VerifyAuthorizationCallCount()).To(Equal(0))
			})
		})

		Context("service fails", func() {
			BeforeEach(func() {
				requirement.VerifyAuthorizationReturns(nil)
			})

			Context("denies authorization", func() {
				It("is not authorized and does not invoke service", func() {
					service.VerifyAuthorizationReturns(errors.New("fake-err1"))

					err := subject.VerifyAuthorization(req, token)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("fake-err1"))
				})
			})
		})

		It("authorizes", func() {
			requirement.VerifyAuthorizationReturns(nil)
			service.VerifyAuthorizationReturns(nil)

			err := subject.VerifyAuthorization(req, token)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
