package config

type Config struct {
	Glob     string            `yaml:"glob,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty"`

	Paths []PathConfig `yaml:"-"`
}

type PathConfig struct {
	Name   string
	Path   string
	Size   int64
	Digest PathDigestConfig
}

type PathDigestConfig struct {
	SHA1   string
	SHA256 string
	SHA512 string
}
