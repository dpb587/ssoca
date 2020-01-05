package onc

// https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md#networkconfiguration-type
type NetworkConfiguration struct {
	GUID     string `json:"GUID,omitempty"`
	Name string `json:"Name,omitempty"`
	VPN      VPN    `json:"VPN,omitempty"`
	Type     string `json:"Type,omitempty"`
	Priority int    `json:"Priority,omitempty"`
}
