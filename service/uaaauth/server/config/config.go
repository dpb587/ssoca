package config

// Config settings for UAA Authentication.
type Config struct {
	// The address of the UAA server. This is used by clients and does not need to be accessible by the server.
	URL string `yaml:"url"`

	// The CA Certificate which the UAA server is secured by (in PEM format).
	CACertificate string `yaml:"ca_certificate,omitempty"`

	// The JWT public key of the UAA server for verifying signed tokens.
	PublicKey string `yaml:"public_key"`

	// The Client ID that ssoca clients should use to connect to UAA. Reminder: this is public to ssoca-client before authentication!
	ClientID string `yaml:"client_id"`

	// The Client Secret that ssoca clients should use to connect to UAA. Reminder: this is public to ssoca-client before authentication!
	ClientSecret string `yaml:"client_secret"`

	// Optionally override which prompts from the UAA server are given to the user.
	Prompts []string `yaml:"prompts"`
}
