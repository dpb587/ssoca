package config

type Config struct {
	Banner        string            `yaml:"banner"`
	Metadata      map[string]string `yaml:"metadata"`
	Name          string            `yaml:"name"`
	Title         string            `yaml:"title"`
	UpdateService string            `yaml:"update_service"`
	URL           string            `yaml:"url"`
}
