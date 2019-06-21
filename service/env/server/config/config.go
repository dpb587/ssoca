package config

type Config struct {
	Banner             string            `yaml:"banner"`
	DefaultAuthService string            `yaml:"default_auth_service"`
	Metadata           map[string]string `yaml:"metadata"`
	Name               string            `yaml:"name"`
	Title              string            `yaml:"title"`
	UpdateService      string            `yaml:"update_service"`
	URL                string            `yaml:"url"`

	SupportOlderClients bool `yaml:"support_older_clients"` // TODO deprecate
}
