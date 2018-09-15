package onc

// https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md#vpn-type
type VPN struct {
	AutoConnect bool    `json:"AutoConnect,omitempty"`
	Host        string  `json:"Host,omitempty"`
	OpenVPN     OpenVPN `json:"OpenVPN,omitempty"`
	Type        string  `json:"Type,omitempty"`
}
