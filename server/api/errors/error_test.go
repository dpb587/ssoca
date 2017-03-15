package errors_test

import (
	"errors"

	. "github.com/dpb587/ssoca/server/api/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	Describe("NewError", func() {
		It("wraps the error", func() {
			err := NewError(errors.New("fake-err"), 123, "fake-message")

			Expect(err.Error()).To(Equal("fake-err"))
			Expect(err.Status).To(Equal(123))
			Expect(err.PublicError).To(Equal("fake-message"))
		})

		It("defaults public message when blank", func() {
			err := NewError(errors.New("fake-err"), 404, "")

			Expect(err.Error()).To(Equal("fake-err"))
			Expect(err.Status).To(Equal(404))
			Expect(err.PublicError).To(Equal("Not Found"))
		})

		It("gives precedence to existing api errors", func() {
			err := NewError(NewError(errors.New("inner-err"), 501, "inner-message"), 401, "outer-message")

			Expect(err.Error()).To(Equal("inner-err"))
			Expect(err.Status).To(Equal(501))
			Expect(err.PublicError).To(Equal("inner-message"))
		})
	})

	Describe("WrapError", func() {
		It("wraps regular errors", func() {
			err := WrapError(errors.New("fake-inner"), "fake-outer")

			Expect(err.Error()).To(Equal("fake-outer: fake-inner"))

			Expect(err).ToNot(BeAssignableToTypeOf(Error{}))
		})

		It("wraps api errors", func() {
			err := WrapError(NewError(errors.New("fake-inner"), 123, "fake-message"), "fake-outer")

			Expect(err.Error()).To(Equal("fake-outer: fake-inner"))

			apiError, ok := err.(Error)
			Expect(ok).To(BeTrue())

			Expect(apiError.Status).To(Equal(123))
			Expect(apiError.PublicError).To(Equal("fake-message"))
		})
	})

	Describe("WrapErrorf", func() {
		It("wraps regular errors", func() {
			err := WrapErrorf(errors.New("fake-inner"), "fake-outer-%s", "formatted")

			Expect(err.Error()).To(Equal("fake-outer-formatted: fake-inner"))

			Expect(err).ToNot(BeAssignableToTypeOf(Error{}))
		})

		It("wraps api errors", func() {
			err := WrapErrorf(NewError(errors.New("fake-inner"), 123, "fake-message"), "fake-outer-%s", "formatted")

			Expect(err.Error()).To(Equal("fake-outer-formatted: fake-inner"))

			apiError, ok := err.(Error)
			Expect(ok).To(BeTrue())

			Expect(apiError.Status).To(Equal(123))
			Expect(apiError.PublicError).To(Equal("fake-message"))
		})
	})
})
