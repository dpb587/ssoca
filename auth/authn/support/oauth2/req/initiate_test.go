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
		It("works", func() {
			subject = Initiate{}

			Expect(subject.Route()).To(Equal("initiate"))
		})
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
					RedirectURL: "https://localhost/somewhere",
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

			stateCookieSplit := strings.SplitN(stateCookie[1], "; ", 2)
			Expect(stateCookieSplit).To(HaveLen(2))
			Expect(stateCookieSplit[1]).To(Equal("Path=/auth/; Domain=localhost; Secure"))

			location, err := url.Parse(res.Header().Get("Location"))

			Expect(err).ToNot(HaveOccurred())
			Expect(location.Host).To(Equal("oauth.example.com"))
			Expect(location.Path).To(Equal("/auth"))
			Expect(location.Scheme).To(Equal("https"))
			Expect(location.Query().Get("client_id")).To(Equal("client-id"))
			Expect(location.Query().Get("response_type")).To(Equal("code"))
			Expect(location.Query().Get("state")).To(Equal(stateCookieSplit[0]))
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
				Expect(portCookie[1]).To(Equal("12345; Path=/auth/; Domain=localhost; Secure"))
			})
		})

		Context("insecure cookies can be used", func() {
			It("sets the cookie", func() {
				subject.Config.RedirectURL = "http://localhost/somewhere"

				err := subject.Execute(req.Request{
					RawRequest:  httptest.NewRequest("GET", "http://localhost/auth/initiate?client_port=12345", nil),
					RawResponse: res,
				})

				Expect(err).ToNot(HaveOccurred())

				Expect(res.Header()["Set-Cookie"]).To(HaveLen(2))

				portCookie := strings.SplitN(res.Header()["Set-Cookie"][1], "=", 2)
				Expect(portCookie).To(HaveLen(2))
				Expect(portCookie[1]).To(Equal("12345; Path=/auth/; Domain=localhost"))
			})
		})

		Context("an aliased host is accessed", func() {
			It("redirects to the correct server", func() {
				subject.Config.RedirectURL = "https://elsewhere.com:54321/somewhere"

				err := subject.Execute(req.Request{
					RawRequest:  httptest.NewRequest("GET", "http://localhost:12345/auth/initiate?client_port=12345", nil),
					RawResponse: res,
				})

				Expect(err).ToNot(HaveOccurred())

				Expect(res.Header()["Set-Cookie"]).To(HaveLen(0))

				Expect(res.Header()["Location"]).To(HaveLen(1))
				Expect(res.Header()["Location"][0]).To(Equal("https://elsewhere.com:54321/auth/initiate?client_port=12345"))
			})
		})
	})
})
