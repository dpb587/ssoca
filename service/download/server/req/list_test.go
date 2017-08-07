package server_test

import (
	"net/http/httptest"

	"github.com/dpb587/ssoca/server/service/req"
	svcconfig "github.com/dpb587/ssoca/service/download/config"
	. "github.com/dpb587/ssoca/service/download/server/req"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	Describe("Execute", func() {
		Context("with files", func() {
			var subject List
			var res httptest.ResponseRecorder

			BeforeEach(func() {
				subject = List{
					Paths: []svcconfig.PathConfig{
						{
							Name: "test1",
							Size: 1234,
							Digest: svcconfig.PathDigestConfig{
								SHA1:   "a1b2c3d4",
								SHA256: "a1b2c3d4a1b2c3d4",
								SHA512: "a1b2c3d4a1b2c3d4a1b2c3d4",
							},
						},
						{
							Name: "test2",
							Size: 5678,
							Digest: svcconfig.PathDigestConfig{
								SHA1:   "e5f6g7h8",
								SHA256: "e5f6g7h8e5f6g7h8",
								SHA512: "e5f6g7h8e5f6g7h8e5f6g7h8",
							},
						},
					},
				}

				res = *httptest.NewRecorder()
			})

			It("enumerates all files and properties", func() {
				err := subject.Execute(req.Request{RawResponse: &res})

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Body.String()).To(Equal(`{
  "files": [
    {
      "name": "test1",
      "size": 1234,
      "digest": {
        "sha1": "a1b2c3d4",
        "sha256": "a1b2c3d4a1b2c3d4",
        "sha512": "a1b2c3d4a1b2c3d4a1b2c3d4"
      }
    },
    {
      "name": "test2",
      "size": 5678,
      "digest": {
        "sha1": "e5f6g7h8",
        "sha256": "e5f6g7h8e5f6g7h8",
        "sha512": "e5f6g7h8e5f6g7h8e5f6g7h8"
      }
    }
  ]
}
`))
			})
		})
	})
})
