package config

type Config struct {
	Path string `yaml:"path,omitempty"`

	AbsPath string `yaml:"-"`
}
