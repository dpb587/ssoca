package dynamicvalue_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/server/service/dynamicvalue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateValue", func() {
	var req *http.Request
	var token *auth.Token

	BeforeEach(func() {
		req = httptest.NewRequest("GET", "/", nil)
		user := "user@fake"
		token = &auth.Token{
			ID: "id@fake",
			Groups: []string{
				"scope1",
				"scope2",
			},
			Attributes: map[auth.TokenAttribute]*string{
				auth.TokenUsernameAttribute: &user,
			},
		}
	})

	Describe("Evaluate", func() {
		Context("simple string", func() {
			It("works", func() {
				subject, err := CreateTemplateValue("hello")

				Expect(err).ToNot(HaveOccurred())
				Expect(subject.Evaluate(req, token)).To(Equal("hello"))
			})
		})

		Context("request info", func() {
			It("works", func() {
				subject, err := CreateTemplateValue("{{ .Request.Method }}")

				Expect(err).ToNot(HaveOccurred())
				Expect(subject.Evaluate(req, token)).To(Equal("GET"))
			})
		})

		Context("token info", func() {
			Context("properties", func() {
				It("works", func() {
					subject, err := CreateTemplateValue(`{{ index ( split .Token.ID "@" ) 0 }}`)

					Expect(err).ToNot(HaveOccurred())
					Expect(subject.Evaluate(req, token)).To(Equal("id"))
				})
			})

			Context("funds", func() {
				It("works", func() {
					subject, err := CreateTemplateValue(`{{ index ( split .Token.Username "@" ) 0 }}`)

					Expect(err).ToNot(HaveOccurred())
					Expect(subject.Evaluate(req, token)).To(Equal("user"))
				})
			})
		})
	})
})
