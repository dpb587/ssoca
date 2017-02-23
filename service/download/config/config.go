package config

type Config struct {
	Glob string `yaml:"glob,omitempty"`

	Paths []PathConfig `yaml:"-"`
}

type PathConfig struct {
	Name   string
	Path   string
	Size   int64
	Digest string
}
