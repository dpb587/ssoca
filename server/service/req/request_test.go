package req_test

import (
	"net/http/httptest"
	"strings"

	. "github.com/dpb587/ssoca/server/service/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type payloadTest struct {
	Test bool `json:"test"`
}

var _ = Describe("Request", func() {
	var subject Request

	Describe("ReadPayload", func() {
		It("unmarshals json", func() {
			subject = Request{
				RawRequest: httptest.NewRequest("POST", "http://localhost/somewhere", strings.NewReader(`{"test":true}`)),
			}

			payload := payloadTest{}

			err := subject.ReadPayload(&payload)

			Expect(err).ToNot(HaveOccurred())
			Expect(payload.Test).To(BeTrue())
		})

		It("captures unmarshal errors", func() {
			subject = Request{
				RawRequest: httptest.NewRequest("POST", "http://localhost/somewhere", strings.NewReader(`{"tes`)),
			}

			err := subject.ReadPayload(struct{}{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unmarshaling request payload"))
		})
	})

	Describe("WritePayload", func() {
		It("marshals json", func() {
			response := httptest.NewRecorder()
			subject = Request{
				RawResponse: response,
			}

			payload := payloadTest{
				Test: true,
			}

			err := subject.WritePayload(payload)

			Expect(err).ToNot(HaveOccurred())
			Expect(response.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(response.Body.String()).To(Equal(`{
  "test": true
}
`))
		})
	})
})
