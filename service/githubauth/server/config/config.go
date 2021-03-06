package config

import (
	"time"

	oauth2config "github.com/dpb587/ssoca/auth/authn/oauth2/server/config"
)

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	AuthURL  string `yaml:"auth_url"`
	TokenURL string `yaml:"token_url"`

	SuccessRedirectURL string `yaml:"success_redirect_url"`
	FailureRedirectURL string `yaml:"failure_redirect_url"`

	JWT JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	PrivateKey   *oauth2config.PrivateKey `yaml:"private_key"`
	Validity     time.Duration            `yaml:"validity"`
	ValidityPast time.Duration            `yaml:"validity_past"`
}

func (c *Config) ApplyDefaults() {
	if c.AuthURL == "" {
		c.AuthURL = "https://github.com/login/oauth/authorize"
	}

	if c.TokenURL == "" {
		c.TokenURL = "https://github.com/login/oauth/access_token"
	}

	c.JWT.ApplyDefaults()
}

func (c *Config) ApplyRedirectDefaults(success, failure string) {
	if c.SuccessRedirectURL == "" {
		c.SuccessRedirectURL = success
	}

	if c.FailureRedirectURL == "" {
		c.FailureRedirectURL = failure
	}
}

func (c *JWTConfig) ApplyDefaults() {
	if c.Validity == 0 {
		v := 24 * time.Hour
		c.Validity = v
	}

	if c.ValidityPast == 0 {
		v := 2 * time.Second
		c.ValidityPast = v
	}
}
