package req_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	. "github.com/dpb587/ssoca/auth/authn/support/oauth2/req"
	"github.com/dpb587/ssoca/server/service/req"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initiate", func() {
	var subject Initiate

	Describe("Route", func() {
		subject = Initiate{}

		Expect(subject.Route()).To(Equal("initiate"))
	})

	Describe("Execute", func() {
		var res http.ResponseWriter

		BeforeEach(func() {
			subject = Initiate{
				Config: oauth2.Config{
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					Endpoint: oauth2.Endpoint{
						AuthURL: "https://oauth.example.com/auth",
					},
				},
			}

			res = httptest.NewRecorder()
		})

		It("works", func() {
			err := subject.Execute(req.Request{
				RawRequest:  httptest.NewRequest("GET", "http://localhost/auth/initiate", nil),
				RawResponse: res,
			})

			Expect(err).ToNot(HaveOccurred())

			Expect(res.Header()["Set-Cookie"]).To(HaveLen(1))

			stateCookie := strings.SplitN(res.Header()["Set-Cookie"][0], "=", 2)
			Expect(stateCookie).To(HaveLen(2))

			location, err := url.Parse(res.Header().Get("Location"))

			Expect(err).ToNot(HaveOccurred())
			Expect(location.Host).To(Equal("oauth.example.com"))
			Expect(location.Path).To(Equal("/auth"))
			Expect(location.Scheme).To(Equal("https"))
			Expect(location.Query().Get("client_id")).To(Equal("client-id"))
			Expect(location.Query().Get("response_type")).To(Equal("code"))
			Expect(location.Query().Get("state")).To(Equal(stateCookie[1]))
		})

		Context("client port passed", func() {
			It("sets the cookie", func() {
				err := subject.Execute(req.Request{
					RawRequest:  httptest.NewRequest("GET", "http://localhost/auth/initiate?client_port=12345", nil),
					RawResponse: res,
				})

				Expect(err).ToNot(HaveOccurred())

				Expect(res.Header()["Set-Cookie"]).To(HaveLen(2))

				portCookie := strings.SplitN(res.Header()["Set-Cookie"][1], "=", 2)
				Expect(portCookie).To(HaveLen(2))
				Expect(portCookie[1]).To(Equal("12345"))
			})
		})
	})
})
