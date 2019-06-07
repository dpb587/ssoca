package config

import (
	"path"
	"strings"
	"time"
)

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	AuthURL  string `yaml:"auth_url"`
	TokenURL string `yaml:"token_url"`

	Scopes ScopesConfig `yaml:"scopes"`

	JWT JWTConfig `yaml:"jwt"`

	RedirectSuccess string `yaml:"redirect_success"`
	RedirectFailure string `yaml:"redirect_failure"`
}

type ScopesConfig struct {
	CloudProject *ScopesCloudProjectConfig `yaml:"cloud_project"`
}

type ScopesCloudProjectConfig struct {
	Projects []string `yaml:"projects"`
	Roles    []string `yaml:"roles"`
}

type JWTConfig struct {
	PrivateKey   string         `yaml:"private_key"`
	Validity     *time.Duration `yaml:"validity"`
	ValidityPast *time.Duration `yaml:"validity_past"`
}

func (c *Config) ApplyDefaults() {
	if c.AuthURL == "" {
		c.AuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
	}

	if c.TokenURL == "" {
		c.TokenURL = "https://www.googleapis.com/oauth2/v4/token"
	}

	c.JWT.ApplyDefaults()
}

func (c *Config) AbsolutizeRedirects(endpoint string) {
	if strings.HasPrefix(c.RedirectSuccess, "/") {
		c.RedirectSuccess = path.Join(endpoint, c.RedirectSuccess)
	}

	if strings.HasPrefix(c.RedirectFailure, "/") {
		c.RedirectFailure = path.Join(endpoint, c.RedirectFailure)
	}
}

func (c *JWTConfig) ApplyDefaults() {
	if c.Validity == nil {
		v := 24 * time.Hour
		c.Validity = &v
	}

	if c.ValidityPast == nil {
		v := 2 * time.Second
		c.ValidityPast = &v
	}
}
