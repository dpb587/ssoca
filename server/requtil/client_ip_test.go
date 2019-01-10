package requtil_test

import (
	"net"
	"net/http"

	. "github.com/dpb587/ssoca/server/requtil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetClientIP", func() {
	mustParseCIDR := func(input string) *net.IPNet {
		_, output, err := net.ParseCIDR(input)
		Expect(err).ToNot(HaveOccurred())

		return output
	}

	It("defaults to request remote address", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "192.0.2.1:19291",
			},
			nil,
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("192.0.2.1"))
	})

	It("trusts a single-level proxy", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "192.168.2.1:19291",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.123"},
				},
			},
			[]*net.IPNet{mustParseCIDR("192.168.0.0/16")},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("203.0.113.123"))
	})

	It("supports IPv6", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "[::1]:1234",
				Header: http.Header{
					"X-Forwarded-For": []string{"2001:db8:85a3:8d3:1319:8a2e:370:7348"},
				},
			},
			[]*net.IPNet{mustParseCIDR("::1/128")},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
	})

	It("trusts a multiple proxies", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "127.0.0.1:19291",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.123, 192.168.2.1"},
				},
			},
			[]*net.IPNet{
				mustParseCIDR("192.168.0.0/16"),
				mustParseCIDR("127.0.0.0/8"),
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("203.0.113.123"))
	})

	It("stops when encountering an untrusted address", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "127.0.0.1:19291",
				Header: http.Header{
					"X-Forwarded-For": []string{"203.0.113.123, 192.168.2.1"},
				},
			},
			[]*net.IPNet{
				mustParseCIDR("127.0.0.0/8"),
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("192.168.2.1"))
	})

	It("stops when encountering an invalid IP", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "127.0.0.1:19291",
				Header: http.Header{
					"X-Forwarded-For": []string{"who.knows.the.where, 192.168.2.1"},
				},
			},
			[]*net.IPNet{
				mustParseCIDR("192.168.0.0/16"),
				mustParseCIDR("127.0.0.0/8"),
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("192.168.2.1"))
	})

	It("uses client IP if proxy is distrusted", func() {
		ip, err := GetClientIP(
			&http.Request{
				RemoteAddr: "127.0.0.1:19291",
				Header: http.Header{
					"X-Forwarded-For": []string{"192.168.2.1"},
				},
			},
			[]*net.IPNet{
				mustParseCIDR("10.0.0.0/24"),
			},
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip.String()).To(Equal("127.0.0.1"))
	})

	It("errors on invalid remote address", func() {
		_, err := GetClientIP(&http.Request{RemoteAddr: "192.hello.world.1"}, nil)
		Expect(err).To(HaveOccurred())
	})
})
