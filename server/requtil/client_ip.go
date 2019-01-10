package requtil

import (
	"net"
	"net/http"
	"strings"
)

type ClientIPGetter func(r *http.Request) (net.IP, error)

func GetClientIPWithoutProxies(r *http.Request) (net.IP, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.IP{}, err
	}

	return net.ParseIP(host), nil
}

func GetClientIP(r *http.Request, trustedProxies []*net.IPNet) (net.IP, error) {
	lastTrustedIP, err := GetClientIPWithoutProxies(r)
	if err != nil {
		return net.IP{}, err
	}

	if len(trustedProxies) == 0 {
		// shortcut and avoid xff parsing attempts
		return lastTrustedIP, nil
	}

	var xff []string

	if header := r.Header.Get("x-forwarded-for"); header != "" {
		xff = strings.SplitN(header, ", ", 64)
	}

	xff = append(xff, lastTrustedIP.String())

	for hostIdx := len(xff) - 1; hostIdx >= 0; hostIdx-- {
		ip := net.ParseIP(xff[hostIdx])
		if ip == nil {
			// bad user data; assume malicious and return last good
			return lastTrustedIP, nil
		}

		if hostIdx == 0 {
			// there's nobody else in the trace; this must be it
			return ip, nil
		}

		// verify they are a trusted proxy
		trusted := false

		for _, trustedProxy := range trustedProxies {
			if !trustedProxy.Contains(ip) {
				continue
			}

			trusted = true

			break
		}

		if !trusted {
			// we don't trust it as a proxy; treat it as the client
			return ip, nil
		}

		lastTrustedIP = ip
	}

	// this should never happen due to explicit earlier returns
	panic("unexpected runtime scenario: trusted proxy verification")
}
