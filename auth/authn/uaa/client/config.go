package client

type AuthConfig struct {
	URL           string `yaml:"url,omitempty"`
	CACertificate string `yaml:"ca_certificate,omitempty"`
	RefreshToken  string `yaml:"refresh_token,omitempty"`
}
