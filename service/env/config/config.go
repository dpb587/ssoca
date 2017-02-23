package config

type Config struct {
	Banner   string            `yaml:"banner"`
	Metadata map[string]string `yaml:"metadata"`
	Name     string            `yaml:"name"`
	Title    string            `yaml:"title"`
	URL      string            `yaml:"url"`
}
