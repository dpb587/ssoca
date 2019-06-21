package server_test

import (
	"errors"

	boshsysfakes "github.com/cloudfoundry/bosh-utils/system/fakes"
	"github.com/dpb587/ssoca/service/file"
	. "github.com/dpb587/ssoca/service/file/server"
	svcconfig "github.com/dpb587/ssoca/service/file/server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Factory", func() {
	var subject ServiceFactory
	var fs boshsysfakes.FakeFileSystem

	BeforeEach(func() {
		fs = *boshsysfakes.NewFakeFileSystem()
		subject = NewServiceFactory(&fs)
	})

	Describe("Type", func() {
		It("works", func() {
			Expect(subject.Type()).To(Equal(file.Type))
		})
	})

	Describe("Create", func() {
		It("works", func() {
			fs.WriteFileString("/testdir/ssoca-server", "server data")
			fs.WriteFileString("/testdir/ssoca-client-one", "one data")
			fs.WriteFileString("/testdir/ssoca-client-two", "two data")
			fs.SetGlob(
				"/testdir/ssoca-client-*",
				[]string{"/testdir/ssoca-client-one", "/testdir/ssoca-client-two"},
			)

			svc, err := subject.Create("test1", map[string]interface{}{
				"glob": "/testdir/ssoca-client-*",
			})

			Expect(err).To(BeNil())
			Expect(svc).ToNot(BeNil())

			downloadSvc, ok := svc.(*Service)
			Expect(ok).To(BeTrue())

			paths := downloadSvc.GetDownloadPaths()

			Expect(paths).To(HaveLen(2))

			Expect(paths[0].Name).To(Equal("ssoca-client-one"))
			Expect(paths[0].Path).To(Equal("/testdir/ssoca-client-one"))
			Expect(paths[0].Size).To(BeEquivalentTo(8))
			Expect(paths[0].Digest).To(Equal(svcconfig.PathDigestConfig{
				SHA1:   "3aa2bfba9635820577e1fec31e8cc3087e2cb003",
				SHA256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				SHA512: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
			}))

			Expect(paths[1].Name).To(Equal("ssoca-client-two"))
			Expect(paths[1].Path).To(Equal("/testdir/ssoca-client-two"))
			Expect(paths[1].Size).To(BeEquivalentTo(8))
			Expect(paths[1].Digest).To(Equal(svcconfig.PathDigestConfig{
				SHA1:   "64acd80fad66f66b398686f0165b1c30edbe3730",
				SHA256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				SHA512: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
			}))
		})

		Context("filesystem errors", func() {
			Describe("glob errors", func() {
				It("errors", func() {
					fs.GlobErr = errors.New("fake-err")

					_, err := subject.Create("test1", map[string]interface{}{
						"glob": "/testdir/ssoca-client-*",
					})

					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(ContainSubstring("globbing"))
					Expect(err.Error()).To(ContainSubstring("fake-err"))
				})
			})

			Describe("fs read errors", func() {
				It("errors", func() {
					fs.WriteFileString("/testdir/ssoca-client-one", "one data")
					fs.SetGlob(
						"/testdir/ssoca-client-*",
						[]string{"/testdir/ssoca-client-one"},
					)
					fs.OpenFileErr = errors.New("fake-err")

					_, err := subject.Create("test1", map[string]interface{}{
						"glob": "/testdir/ssoca-client-*",
					})

					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(ContainSubstring("opening file for digest"))
					Expect(err.Error()).To(ContainSubstring("fake-err"))
				})
			})
		})
	})
})
