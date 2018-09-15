package onc

// https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md
type ONC struct {
	Type                  string                 `json:"Type,omitempty"`
	NetworkConfigurations []NetworkConfiguration `json:"NetworkConfigurations,omitempty"`
	Certificates          []Certificate          `json:"Certificates,omitempty"`
}
