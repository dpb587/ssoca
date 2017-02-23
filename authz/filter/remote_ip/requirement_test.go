package remote_ip_test

import (
	"net"
	"net/http"

	. "github.com/dpb587/ssoca/authz/filter/remote_ip"

	"github.com/dpb587/ssoca/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requirement", func() {
	var token auth.Token
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

	BeforeEach(func() {
		token = auth.NewSimpleToken("test", map[string]interface{}{})
	})

	Describe("IsSatisfied", func() {
		Context("IP matching", func() {
			BeforeEach(func() {
				subject = createRequirement("192.0.2.29/32")
			})

			It("satisfies with IP", func() {
				satisfied, err := subject.IsSatisfied(
					&http.Request{
						RemoteAddr: "192.0.2.29:1234",
					},
					token,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeTrue())
			})

			It("does not match another IP", func() {
				satisfied, err := subject.IsSatisfied(
					&http.Request{
						RemoteAddr: "192.0.2.28:1234",
					},
					token,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeFalse())
			})
		})

		Context("netmask matching", func() {
			BeforeEach(func() {
				subject = createRequirement("192.0.2.0/24")
			})

			It("satisfies with netmask", func() {
				satisfied, err := subject.IsSatisfied(
					&http.Request{
						RemoteAddr: "192.0.2.27:1234",
					},
					token,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeTrue())
			})

			It("does not match another IP", func() {
				satisfied, err := subject.IsSatisfied(
					&http.Request{
						RemoteAddr: "192.0.3.27:1234",
					},
					token,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(satisfied).To(BeFalse())
			})
		})

		Context("eccentricities", func() {
			It("error if request remote address is malformed", func() {
				satisfied, err := subject.IsSatisfied(
					&http.Request{
						RemoteAddr: "impractical",
					},
					token,
				)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Parsing remote address"))
				Expect(satisfied).To(BeFalse())
			})
		})
	})
})
