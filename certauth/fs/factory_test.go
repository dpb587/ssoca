package fs_test

import (
	"github.com/Sirupsen/logrus"
	logrustest "github.com/Sirupsen/logrus/hooks/test"
	. "github.com/dpb587/ssoca/certauth/fs"

	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	var subject Factory

	Describe("Create", func() {
		var fs boshsysfakes.FakeFileSystem
		var logger logrus.FieldLogger

		BeforeEach(func() {
			fs = *boshsysfakes.NewFakeFileSystem()
			logger, _ = logrustest.NewNullLogger()

			subject = NewFactory(&fs, logger)
		})

		It("remarshals configuration", func() {
			provider, err := subject.Create("name1", map[string]interface{}{
				"certificate_path": "/somewhere",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(provider.Name()).To(Equal("name1"))
		})

		Context("invalid yaml", func() {
			It("remarshals configuration", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certificate_path": map[string]interface{}{},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Loading config"))
			})
		})
	})
})
