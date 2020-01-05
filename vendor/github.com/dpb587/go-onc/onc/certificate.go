package onc

// https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md#certificate-type
type Certificate struct {
	GUID      string   `json:"GUID,omitempty"`
	PKCS12    string   `json:"PKCS12,omitempty"`
	Remove    bool     `json:"Remove,omitempty"`
	TrustBits []string `json:"TrustBits,omitempty"`
	Type      string   `json:"Type,omitempty"`
	X509      string   `json:"X509,omitempty"`
}
