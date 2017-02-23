package config

import "time"

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	AuthURL  string `yaml:"auth_url"`
	TokenURL string `yaml:"token_url"`

	JWT JWTConfig `yaml:"jwt"`
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
		c.AuthURL = "https://github.com/login/oauth/authorize"
	}

	if c.TokenURL == "" {
		c.TokenURL = "https://github.com/login/oauth/access_token"
	}
}
