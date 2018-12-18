package client

type AuthConfig struct {
	URL           string `yaml:"url,omitempty"`
	CACertificate string `yaml:"ca_certificate,omitempty"`
	ClientID      string `yaml:"client_id"`
	ClientSecret  string `yaml:"client_secret"`
	RefreshToken  string `yaml:"refresh_token,omitempty"`
}
