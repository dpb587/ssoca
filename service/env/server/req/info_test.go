package req_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/server/service/servicefakes"
	svcapi "github.com/dpb587/ssoca/service/env/api"
	"github.com/dpb587/ssoca/service/env/server/config"
	. "github.com/dpb587/ssoca/service/env/server/req"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	var subject Info
	var token auth.Token
	var loggerContext logrus.Fields
	var res httptest.ResponseRecorder

	BeforeEach(func() {
		loggerContext = logrus.Fields{
			"custom": "fake",
		}

		svc1 := &servicefakes.FakeService{}
		svc1.NameReturns("env")
		svc1.TypeReturns("fake-env-type")
		svc1.VersionReturns("fake-env-version")
		svc1.MetadataReturns(map[string]string{"a": "b"})

		svc2 := &servicefakes.FakeService{}
		svc2.VerifyAuthorizationReturns(errors.New("fake-err"))

		svc3 := &servicefakes.FakeService{}
		svc3.NameReturns("fake-auth-name")
		svc3.TypeReturns("fake-auth-type")
		svc3.VersionReturns("fake-auth-version")
		svc3.MetadataReturns(map[string]string{"c": "d"})

		manager := service.NewDefaultManager()
		manager.Add(svc1)
		manager.Add(svc2)
		manager.Add(svc3)

		subject = Info{
			Config: config.Config{
				Banner:             "fake-banner",
				DefaultAuthService: "fake-auth-name",
				Metadata: map[string]string{
					"fake-metadata-key": "fake-metadata-value",
				},
				Name:          "fake-name",
				Title:         "fake-title",
				UpdateService: "fake-update-service",
				URL:           "https://fake-url:12345",
			},
			Services: manager,
		}

		token = auth.Token{ID: "fake-user"}
		res = *httptest.NewRecorder()
	})

	Describe("Execute", func() {
		It("works", func() {
			req := req.Request{
				RawRequest:    httptest.NewRequest("GET", "https://localhost/info", nil),
				RawResponse:   &res,
				AuthToken:     &token,
				LoggerContext: loggerContext,
			}

			err := subject.Execute(req)
			Expect(err).ToNot(HaveOccurred())

			var resPayload svcapi.InfoResponse

			err = json.Unmarshal(res.Body.Bytes(), &resPayload)
			Expect(err).ToNot(HaveOccurred())

			Expect(resPayload.Auth).To(BeNil())
			Expect(resPayload.Env.Name).To(Equal("fake-name"))
			Expect(resPayload.Version).To(Equal("fake-env-version")) // TODO remove this api field?
			Expect(resPayload.Services).To(HaveLen(1))               // not env; not unauthorized
			Expect(resPayload.Services[0].Type).To(Equal("fake-auth-type"))
			Expect(resPayload.Services[0].Name).To(Equal("fake-auth-name"))
			Expect(resPayload.Services[0].Metadata).To(Equal(map[string]interface{}{"c": "d"}))
		})

		Context("legacy API responses", func() {
			BeforeEach(func() {
				subject.Config.SupportOlderClients = true
			})

			It("populates auth", func() {
				req := req.Request{
					RawRequest:    httptest.NewRequest("GET", "https://localhost/info", nil),
					RawResponse:   &res,
					AuthToken:     &token,
					LoggerContext: loggerContext,
				}

				err := subject.Execute(req)
				Expect(err).ToNot(HaveOccurred())

				var resPayload svcapi.InfoResponse

				err = json.Unmarshal(res.Body.Bytes(), &resPayload)
				Expect(err).ToNot(HaveOccurred())

				Expect(resPayload.Auth).ToNot(BeNil())
				Expect(resPayload.Auth.Type).To(Equal("fake-auth-type"))
				Expect(resPayload.Auth.Name).To(Equal("fake-auth-name"))
				Expect(resPayload.Auth.Metadata).To(Equal(map[string]interface{}{"c": "d"}))
				Expect(resPayload.Services).To(HaveLen(1))
			})
		})
	})
})
