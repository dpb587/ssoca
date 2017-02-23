package server_test

import (
	svcconfig "github.com/dpb587/ssoca/service/download/config"
	. "github.com/dpb587/ssoca/service/download/server/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	Describe("Execute", func() {
		Context("with files", func() {
			var subject List

			BeforeEach(func() {
				subject = List{
					Paths: []svcconfig.PathConfig{
						{
							Name:   "test1",
							Size:   1234,
							Digest: "a1b2c3d4",
						},
						{
							Name:   "test2",
							Size:   5678,
							Digest: "e5f6g7h8",
						},
					},
				}
			})

			It("enumerates all files and properties", func() {
				res, err := subject.Execute()

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Files).To(HaveLen(2))
				Expect(res.Files[0].Digest).To(Equal("a1b2c3d4"))
				Expect(res.Files[0].Name).To(Equal("test1"))
				Expect(res.Files[0].Size).To(BeEquivalentTo(1234))
				Expect(res.Files[1].Digest).To(Equal("e5f6g7h8"))
				Expect(res.Files[1].Name).To(Equal("test2"))
				Expect(res.Files[1].Size).To(BeEquivalentTo(5678))
			})
		})
	})
})
