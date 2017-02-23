package auth_test

import (
	. "github.com/dpb587/ssoca/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SimpleToken", func() {
	Describe("Username", func() {
		token := NewSimpleToken("user1", map[string]interface{}{})

		Expect(token.Username()).To(Equal("user1"))
	})

	Context("attributes", func() {
		var token Token
		var attributes = map[string]interface{}{
			"key1": "val1",
			"key2": false,
			"key3": []string{
				"value1",
				"value2",
			},
			"key4": nil,
		}

		BeforeEach(func() {
			token = NewSimpleToken("user1", attributes)
		})

		Describe("Attributes", func() {
			It("returns", func() {
				Expect(token.Attributes()).To(Equal(attributes))
			})
		})

		Describe("HasAttribute", func() {
			Context("when present", func() {
				It("returns true", func() {
					Expect(token.HasAttribute("key1")).To(BeTrue())
				})
			})

			Context("when value is nil", func() {
				It("returns true", func() {
					Expect(token.HasAttribute("key4")).To(BeTrue())
				})
			})

			Context("when missing", func() {
				It("returns false", func() {
					Expect(token.HasAttribute("keyA")).To(BeFalse())
				})
			})
		})

		Describe("GetAttribute", func() {
			Context("when present", func() {
				It("returns value", func() {
					value, err := token.GetAttribute("key1")

					Expect(err).ToNot(HaveOccurred())
					Expect(value).To(Equal("val1"))
				})
			})

			Context("when missing", func() {
				It("returns error", func() {
					value, err := token.GetAttribute("keyA")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(": keyA"))
					Expect(value).To(BeNil())
				})
			})
		})
	})
})
