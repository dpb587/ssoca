package config

import (
	"time"

	"github.com/dpb587/ssoca/certauth"
)

type Config struct {
	CertAuth certauth.ConfigValue `yaml:"certauth,omitempty"`
	Validity time.Duration        `yaml:"validity,omitempty"`
	Profile  string               `yaml:"profile,omitempty"`
}

func (c *Config) ApplyDefaults() {
	if !c.CertAuth.IsConfigured() {
		err := c.CertAuth.Configure(certauth.DefaultName)
		if err != nil {
			panic(err)
		}
	}
}
