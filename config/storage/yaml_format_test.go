package storage_test

import (
	"errors"

	. "github.com/dpb587/ssoca/config/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type yamlCustomStructTags struct {
	CustomField bool        `yaml:"something"`
	CustomValue interface{} `yaml:"else,omitempty"`
}

type yamlErrorMarshal struct{}

func (yamlErrorMarshal) MarshalText() (text []byte, err error) {
	return []byte{}, errors.New("fake-error")
}

var _ = Describe("YamlFormat", func() {
	var subject YAMLFormat

	BeforeEach(func() {
		subject = YAMLFormat{}
	})

	Describe("Get", func() {
		var get yamlCustomStructTags

		BeforeEach(func() {
			get = yamlCustomStructTags{}
		})

		It("parses tagged struct", func() {
			err := subject.Get(`something: true`, &get)

			Expect(err).ToNot(HaveOccurred())
			Expect(get.CustomField).To(BeTrue())
		})

		It("can error", func() {
			err := subject.Get(`something: "string"`, &get)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unmarshaling"))
		})
	})

	Describe("Put", func() {
		It("parses tagged struct", func() {
			put := yamlCustomStructTags{
				CustomField: true,
			}

			yaml, err := subject.Put("", put)

			Expect(err).ToNot(HaveOccurred())
			Expect(yaml).To(Equal(`something: true
`))
		})

		It("can error", func() {
			put := yamlCustomStructTags{
				CustomValue: yamlErrorMarshal{},
			}

			_, err := subject.Put("", put)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-error"))
			Expect(err.Error()).To(ContainSubstring("Marshaling"))
		})
	})
})
