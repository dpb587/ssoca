package remote_ip_test

import (
	"net"
	"net/http"

	"github.com/dpb587/ssoca/auth/authz"
	. "github.com/dpb587/ssoca/auth/authz/filter/remote_ip"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requirement", func() {
	var subject Requirement

	createRequirement := func(netmask string) Requirement {
		_, cidr, err := net.ParseCIDR(netmask)
		if err != nil {
			panic(err)
		}

		return Requirement{
			Within: *cidr,
		}
	}

	Describe("VerifyAuthorization", func() {
		Context("IP matching", func() {
			BeforeEach(func() {
				subject = createRequirement("192.0.2.29/32")
			})

			It("satisfies with IP", func() {
				err := subject.VerifyAuthorization(
					&http.Request{
						RemoteAddr: "192.0.2.29:1234",
					},
					nil,
				)

				Expect(err).ToNot(HaveOccurred())
			})

			It("does not match another IP", func() {
				err := subject.VerifyAuthorization(
					&http.Request{
						RemoteAddr: "192.0.2.28:1234",
					},
					nil,
				)

				Expect(err).To(HaveOccurred())

				aerr, ok := err.(authz.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("Remote IP is not allowed"))
			})
		})

		Context("netmask matching", func() {
			BeforeEach(func() {
				subject = createRequirement("192.0.2.0/24")
			})

			It("satisfies with netmask", func() {
				err := subject.VerifyAuthorization(
					&http.Request{
						RemoteAddr: "192.0.2.27:1234",
					},
					nil,
				)

				Expect(err).ToNot(HaveOccurred())
			})

			It("does not match another IP", func() {
				err := subject.VerifyAuthorization(
					&http.Request{
						RemoteAddr: "192.0.3.27:1234",
					},
					nil,
				)

				Expect(err).To(HaveOccurred())

				aerr, ok := err.(authz.Error)
				Expect(ok).To(BeTrue())
				Expect(aerr.Error()).To(Equal("Remote IP is not allowed"))
			})
		})

		Context("eccentricities", func() {
			It("error if request remote address is malformed", func() {
				err := subject.VerifyAuthorization(
					&http.Request{
						RemoteAddr: "impractical",
					},
					nil,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing remote address"))
			})
		})
	})
})
