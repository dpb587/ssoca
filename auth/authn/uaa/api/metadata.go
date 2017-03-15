package api

type Metadata struct {
	URL           string `json:"url"`
	CACertificate string `json:"ca_certificate,omitempty"`
}
