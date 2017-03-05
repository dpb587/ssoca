package dynamicvalue_test

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	. "github.com/dpb587/ssoca/server/service/dynamicvalue"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigValue", func() {
	var subject ConfigValue

	Describe("UnmarshalYAML", func() {
		BeforeEach(func() {
			subject = NewConfigValue(DefaultFactory{})
		})

		It("unmarshals", func() {
			err := yaml.Unmarshal([]byte("test-config"), &subject)

			Expect(err).ToNot(HaveOccurred())

			result, err := subject.Evaluate(&http.Request{}, &auth.Token{})

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("test-config"))
		})

		It("propagates factory errors", func() {
			err := yaml.Unmarshal([]byte("'{{ missing_func }}'"), &subject)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing_func"))
		})
	})
})
