package config

import "time"

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	AuthURL  string `yaml:"auth_url"`
	TokenURL string `yaml:"token_url"`

	Scopes ScopesConfig `yaml:"scopes"`

	JWT JWTConfig `yaml:"jwt"`
}

type ScopesConfig struct {
	CloudProject *ScopesCloudProjectConfig `yaml:"cloud_project"`
}

type ScopesCloudProjectConfig struct {
	Projects []string `yaml:"projects"`
	Roles    []string `yaml:"roles"`
}

type JWTConfig struct {
	PrivateKey         string `yaml:"private_key"`
	ValidityString     string `yaml:"validity"`
	ValidityPastString string `yaml:"validity_past"`

	Validity     time.Duration `yaml:"-"`
	ValidityPast time.Duration `yaml:"-"`
}

func (c *Config) ApplyDefaults() {
	if c.JWT.ValidityString == "" {
		c.JWT.ValidityString = "24h"
	}

	if c.JWT.ValidityPastString == "" {
		c.JWT.ValidityPastString = "5s"
	}

	if c.AuthURL == "" {
		c.AuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
	}

	if c.TokenURL == "" {
		c.TokenURL = "https://www.googleapis.com/oauth2/v4/token"
	}
}
