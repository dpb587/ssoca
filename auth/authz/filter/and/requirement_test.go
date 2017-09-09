package and_test

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/auth/authz/filter/and"

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

			It("satisfies", func() {
				err := subject.VerifyAuthorization(&request, token)

				Expect(err).ToNot(HaveOccurred())
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
			It("does not satisfy", func() {
				subject = Requirement{
					Requirements: []filter.Requirement{
						satisfyRequirement,
						dissatisfyRequirement,
					},
				}

				err := subject.VerifyAuthorization(&request, token)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("faked error"))
			})

			It("stops evaluating early", func() {
				subject = Requirement{
					Requirements: []filter.Requirement{
						satisfyRequirement,
						dissatisfyRequirement,
						satisfyRequirement,
					},
				}

				err := subject.VerifyAuthorization(&request, token)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("faked error"))
				Expect(satisfyRequirement.VerifyAuthorizationCallCount()).To(Equal(1))
			})
		})
	})
})
