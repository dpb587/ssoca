package server_test

import (
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/certauth/certauthfakes"
	"github.com/dpb587/ssoca/server/service/dynamicvalue/dynamicvaluefakes"
	. "github.com/dpb587/ssoca/service/ssh/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceFactory", func() {
	var subject ServiceFactory

	Describe("Type", func() {
		It("returns", func() {
			Expect(subject.Type()).To(Equal("ssh"))
		})
	})

	Describe("Create", func() {
		var caManager certauth.Manager
		var fakedynamicvalue *dynamicvaluefakes.FakeFactory

		BeforeEach(func() {
			caManager = certauth.NewDefaultManager()

			certauth := certauthfakes.FakeProvider{}
			certauth.NameReturns("default")
			caManager.Add(&certauth)

			fakedynamicvalue = &dynamicvaluefakes.FakeFactory{}

			subject = NewServiceFactory(fakedynamicvalue, caManager)
		})

		It("remarshals configuration", func() {
			provider, err := subject.Create("name1", map[string]interface{}{
				"extensions": []string{"permit-user-rc"},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(provider.Name()).To(Equal("name1"))
		})

		Context("invalid certauth", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"certauth": "unknown",
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("getting certificate authority"))
			})
		})

		Context("invalid validity duration", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"validity": "525,600 minutes",
				})

				Expect(err).To(HaveOccurred())
			})
		})

		Context("invalid yaml", func() {
			It("errors", func() {
				_, err := subject.Create("name1", map[string]interface{}{
					"extensions": map[string]interface{}{},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("loading config"))
			})
		})

		Describe("target public key", func() {
			Context("bad PEM block", func() {
				It("errors", func() {
					_, err := subject.Create("name1", map[string]interface{}{
						"target": map[string]string{
							"public_key": `-----BEGIN PUBLIC KEY-----
`,
						},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to parse public key"))
				})
			})

			Context("bad PEM block", func() {
				It("errors", func() {
					_, err := subject.Create("name1", map[string]interface{}{
						"target": map[string]string{
							"public_key": `-----BEGIN PUBLIC KEY-----
hQIDAQAB
-----END PUBLIC KEY-----
`,
						},
					})

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("parsing public key"))
				})
			})

			Context("PEM block", func() {
				It("remarshals configuration", func() {
					provider, err := subject.Create("name1", map[string]interface{}{
						"target": map[string]string{
							"public_key": `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0w9bqwoQkFFKVNSQFvod
EynIGjZK7pzQF876f7gr0WxxZ2QKU8T/q3YvsR9h2RifgJucPXY/lvnxomkwtbnM
rPKf7VGO2M4RaNbNJP38AnJJoz7+A5a8qLPuXIv52PnG4K6GYfkcI+w9A3pIdE9o
vNz3/yg+8lbnngJE7H8OYhmusUHRBdsOjk0Z7CYaOUrQHcTqFXXf1r9JgZTB1T1r
1TcF67XL84ZqjKEhNO2bi0/CAn9j4e/KMGf4fZl6+6ONRcVGUXrq7rBbm12657cW
TcvnU1cVWRlKGQSBs/8gsy8zvOTsz2C4Uyf1w3adVNKY08U+1qtLW+i+ZBUc/m96
hQIDAQAB
-----END PUBLIC KEY-----
`,
						},
					})

					Expect(err).ToNot(HaveOccurred())
					Expect(provider.Name()).To(Equal("name1"))
				})
			})
		})
	})
})
