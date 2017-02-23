package storage_test

import (
	"errors"

	. "github.com/dpb587/ssoca/config/storage"

	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FormattedFs", func() {
	var subject FormattedFS
	var fs boshsysfakes.FakeFileSystem

	BeforeEach(func() {
		fs = *boshsysfakes.NewFakeFileSystem()

		subject = NewFormattedFS(&fs, YAMLFormat{})
	})

	Describe("Get", func() {
		var get string

		BeforeEach(func() {
			get = ""

			fs.WriteFileString("/fake/data", `"test"`)
		})

		It("works", func() {
			err := subject.Get("/fake/data", &get)

			Expect(err).ToNot(HaveOccurred())
			Expect(get).To(Equal("test"))
		})

		// @todo debate
		It("works if file does not exist", func() {
			err := subject.Get("/fake/nonexistant", &get)

			Expect(err).ToNot(HaveOccurred())
			Expect(get).To(Equal(""))
		})

		Context("expanding paths", func() {
			It("works", func() {
				fs.ExpandPathExpanded = "/fake/data"

				err := subject.Get("~/data", &get)

				Expect(err).ToNot(HaveOccurred())
				Expect(get).To(Equal("test"))
			})

			Context("when it fails", func() {
				It("errors", func() {
					fs.ExpandPathErr = errors.New("fake-error")

					err := subject.Get("~/data", &get)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-error"))
				})
			})
		})

		Context("file errors", func() {
			It("errors", func() {
				fs.ReadFileError = errors.New("fake-error")

				err := subject.Get("/fake/data", &get)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Reading config file '/fake/data'"))
			})
		})

		Context("parser errors", func() {
			It("errors", func() {
				fs.WriteFileString("/fake/data", `"test`)

				err := subject.Get("/fake/data", &get)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing config"))
			})
		})
	})

	Describe("Put", func() {
		It("works", func() {
			path, err := subject.Put("/fake/data", "colon: something")

			Expect(err).ToNot(HaveOccurred())
			Expect(path).To(Equal("/fake/data"))

			put, err := fs.ReadFileString("/fake/data")

			Expect(err).ToNot(HaveOccurred())
			Expect(put).To(Equal("'colon: something'\n"))
		})

		Context("expanding paths", func() {
			It("works", func() {
				fs.ExpandPathExpanded = "/fake/data"

				_, err := subject.Put("~/data", "test")

				Expect(err).ToNot(HaveOccurred())

				put, err := fs.ReadFileString("/fake/data")

				Expect(err).ToNot(HaveOccurred())
				Expect(put).To(Equal("test\n"))
			})

			Context("when it fails", func() {
				It("errors", func() {
					fs.ExpandPathErr = errors.New("fake-error")

					_, err := subject.Put("~/data", "test")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("fake-error"))
				})
			})
		})

		Context("parser errors", func() {
			It("errors", func() {
				_, err := subject.Put("/fake/data", yamlErrorMarshal{})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-error"))
				Expect(err.Error()).To(ContainSubstring("Serializing config"))
			})
		})

		Context("file errors", func() {
			It("errors", func() {
				fs.WriteFileError = errors.New("fake-error")

				_, err := subject.Put("/fake/data", "test")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Writing config file '/fake/data'"))
			})
		})
	})
})
