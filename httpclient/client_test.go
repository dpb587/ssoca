package httpclient_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/dpb587/ssoca/httpclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type jsonResponse struct {
	Message string `json:"message"`
	Works   bool   `json:"works"`
}

type badMarshal struct{}

func (badMarshal) MarshalJSON() ([]byte, error) {
	return nil, errors.New("fake-marshal-err")
}

type mockTransport struct {
	rt func(req *http.Request) (resp *http.Response, err error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return t.rt(req)
}

var _ = Describe("Client", func() {
	var subject *Client

	BeforeEach(func() {
		subject = NewClient("https://example.com/subpath", nil)
	})

	Describe("ExpandURI", func() {
		It("expands relative URIs", func() {
			Expect(subject.ExpandURI("/elsewhere")).To(Equal("https://example.com/subpath/elsewhere"))
		})

		It("expands relative URIs", func() {
			Expect(subject.ExpandURI("https://abs.example.com/elsewhere")).To(Equal("https://abs.example.com/elsewhere"))
		})
	})

	Describe("APIGet", func() {
		It("unmarshals successful responses", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(`{"message":"body","works":true}`)),
						}, nil
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			response := &jsonResponse{}

			err := subject.APIGet("/test1", &response)
			Expect(err).ToNot(HaveOccurred())

			Expect(response.Message).To(Equal("body"))
			Expect(response.Works).To(BeTrue())
		})

		It("errors on tcp failures", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return nil, errors.New("fake-err")
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			err := subject.APIGet("/test1", &jsonResponse{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
			Expect(err.Error()).To(ContainSubstring("Executing request"))
		})

		It("errors on non-successful responses", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return &http.Response{
							StatusCode: 403,
							Body:       ioutil.NopCloser(strings.NewReader(`Forbidden`)),
						}, nil
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			err := subject.APIGet("/test1", &jsonResponse{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("HTTP 403: Forbidden"))
		})
	})

	Describe("APIPost", func() {
		It("marshals requests and unmarshals successful responses", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						bodyBytes, _ := ioutil.ReadAll(r.Body)

						Expect(string(bodyBytes)).To(Equal(`{"reqtest1":"value1","reqtest2":true}`))

						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(`{"message":"body","works":true}`)),
						}, nil
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			response := &jsonResponse{}

			err := subject.APIPost("/test1", &response, &map[string]interface{}{"reqtest1": "value1", "reqtest2": true})
			Expect(err).ToNot(HaveOccurred())

			Expect(response.Message).To(Equal("body"))
			Expect(response.Works).To(BeTrue())
		})

		It("errors cleanly with marshal errors", func() {
			err := subject.APIPost("/test1", &jsonResponse{}, badMarshal{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-marshal-err"))
			Expect(err.Error()).To(ContainSubstring("Marshaling request body"))
		})

		It("errors on tcp failures", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return nil, errors.New("fake-err")
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			err := subject.APIPost("/test1", &jsonResponse{}, struct{}{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
			Expect(err.Error()).To(ContainSubstring("Executing request"))
		})

		It("errors on non-successful responses", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return &http.Response{
							StatusCode: 403,
							Body:       ioutil.NopCloser(strings.NewReader(`Forbidden`)),
						}, nil
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			err := subject.APIPost("/test1", &jsonResponse{}, struct{}{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("HTTP 403: Forbidden"))
		})

		It("errors on bad server json", func() {
			subject.Client.Transport = &mockTransport{
				rt: func(r *http.Request) (w *http.Response, err error) {
					switch r.URL.String() {
					case "https://example.com/subpath/test1":
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(strings.NewReader(`!json`)),
						}, nil
					}

					Fail(fmt.Sprintf("unexpected request: %s", r.URL.String()))

					return &http.Response{}, nil
				},
			}

			err := subject.APIPost("/test1", &jsonResponse{}, struct{}{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unmarshaling response body"))
		})
	})
})
