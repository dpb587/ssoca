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
	PrivateKey   string        `yaml:"private_key"`
	Validity     time.Duration `yaml:"validity"`
	ValidityPast time.Duration `yaml:"validity_past"`
}
