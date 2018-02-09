package or_test

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz"
	. "github.com/dpb587/ssoca/auth/authz/filter/or"

	"github.com/dpb587/ssoca/auth/authz/filter"
	"github.com/dpb587/ssoca/auth/authz/filter/filterfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requirement", func() {
	var request http.Request
	var token *auth.Token
	var subject Requirement
	var satisfyRequirement, dissatisfyRequirement *filterfakes.FakeRequirement

	BeforeEach(func() {
		request = http.Request{}
		token = &auth.Token{ID: "test"}

		satisfyRequirement = &filterfakes.FakeRequirement{}

		dissatisfyRequirement = &filterfakes.FakeRequirement{}
		dissatisfyRequirement.VerifyAuthorizationReturns(errors.New("faked error"))
	})

	Describe("VerifyAuthorization", func() {
		Context("without filters", func() {
			BeforeEach(func() {
				subject = Requirement{
					Requirements: []filter.Requirement{},
				}
			})

			It("does not satisfy", func() {
				err := subject.VerifyAuthorization(&request, token)
				Expect(err).To(HaveOccurred())

				err, ok := err.(authz.Error)
				Expect(ok).To(BeTrue())
				Expect(err.Error()).To(Equal("No filters authorized access"))
			})
		})

		Context("with satisfying filters", func() {
			BeforeEach(func() {
				subject = Requirement{
					Requirements: []filter.Requirement{
						satisfyRequirement,
						satisfyRequirement,
					},
				}
			})

			It("satisfies", func() {
				err := subject.VerifyAuthorization(&request, token)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with some dissatisfying filters", func() {
			It("satisfies", func() {
				subject = Requirement{
					Requirements: []filter.Requirement{
						dissatisfyRequirement,
						satisfyRequirement,
					},
				}

				err := subject.VerifyAuthorization(&request, token)

				Expect(err).ToNot(HaveOccurred())
			})

			It("stops evaluating early", func() {
				subject = Requirement{
					Requirements: []filter.Requirement{
						satisfyRequirement,
						dissatisfyRequirement,
					},
				}

				err := subject.VerifyAuthorization(&request, token)

				Expect(err).ToNot(HaveOccurred())
				Expect(dissatisfyRequirement.VerifyAuthorizationCallCount()).To(Equal(0))
			})
		})
	})
})
