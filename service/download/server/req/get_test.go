package server_test

import (
	"errors"

	apierr "github.com/dpb587/ssoca/server/api/errors"
	svcconfig "github.com/dpb587/ssoca/service/download/server/config"
	. "github.com/dpb587/ssoca/service/download/server/req"

	"net/http/httptest"

	"github.com/dpb587/ssoca/server/service/req"

	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get", func() {
	var subject Get
	var fs boshsysfakes.FakeFileSystem

	BeforeEach(func() {
		fs = *boshsysfakes.NewFakeFileSystem()
		fs.WriteFileString("/test1", "test1 data")
		fs.WriteFileString("/test2", "test2 data")

		subject = Get{
			Paths: []svcconfig.PathConfig{
				{
					Name: "test1",
					Path: "/test1",
					Size: 12345,
				},
			},
			FS: &fs,
		}
	})

	Describe("Execute", func() {
		var res httptest.ResponseRecorder

		BeforeEach(func() {
			res = *httptest.NewRecorder()
		})

		It("works", func() {
			err := subject.Execute(
				req.Request{
					RawRequest:  httptest.NewRequest("GET", "https://localhost/file?name=test1", nil),
					RawResponse: &res,
				},
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(res.Header().Get("Content-Disposition")).To(Equal(`attachment; filename="test1"`))
			Expect(res.Header().Get("Content-Length")).To(Equal("12345"))
			Expect(res.Body.String()).To(Equal("test1 data"))
		})

		Context("missing name query", func() {
			It("errors with 404", func() {
				err := subject.Execute(
					req.Request{
						RawRequest:  httptest.NewRequest("GET", "https://localhost/file", nil),
						RawResponse: &res,
					},
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Missing query parameter: name"))

				apiErr, ok := err.(apierr.Error)
				Expect(ok).To(BeTrue())
				Expect(apiErr.Status).To(Equal(404))
			})
		})

		Context("unregistered file request", func() {
			It("errors with 404", func() {
				err := subject.Execute(
					req.Request{
						RawRequest:  httptest.NewRequest("GET", "https://localhost/file?name=nonexistant", nil),
						RawResponse: &res,
					},
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid file name"))

				apiErr, ok := err.(apierr.Error)
				Expect(ok).To(BeTrue())
				Expect(apiErr.Status).To(Equal(404))
			})

			Context("when file exists", func() {
				It("errors with 404", func() {
					err := subject.Execute(
						req.Request{
							RawRequest:  httptest.NewRequest("GET", "https://localhost/file?name=/test2", nil),
							RawResponse: &res,
						},
					)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid file name"))

					apiErr, ok := err.(apierr.Error)
					Expect(ok).To(BeTrue())
					Expect(apiErr.Status).To(Equal(404))
				})
			})
		})

		Context("filesystem errors", func() {
			It("errors", func() {
				fs.OpenFileErr = errors.New("fake-err")

				err := subject.Execute(
					req.Request{
						RawRequest:  httptest.NewRequest("GET", "https://localhost/file?name=test1", nil),
						RawResponse: &res,
					},
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Opening file for reading"))
			})
		})
	})
})
