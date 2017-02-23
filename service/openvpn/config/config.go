package config

import (
	"time"

	"github.com/dpb587/ssoca/certauth"
)

type Config struct {
	CertAuthName   string `yaml:"certauth,omitempty"`
	ValidityString string `yaml:"validity,omitempty"`
	Profile        string `yaml:"profile,omitempty"`

	CertAuth certauth.Provider `yaml:"-"`
	Validity time.Duration     `yaml:"-"`
}

func (c *Config) ApplyDefaults() {
	if c.CertAuthName == "" {
		c.CertAuthName = "default"
	}

	if c.ValidityString == "" {
		c.ValidityString = "2m"
	}
}
