package api

type Metadata struct {
	URL           string   `json:"url"`
	CACertificate string   `json:"ca_certificate,omitempty"`
	ClientID      string   `json:"client_id"`
	ClientSecret  string   `json:"client_secret"`
	Prompts       []string `json:"prompts"`
}
