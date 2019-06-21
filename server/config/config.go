package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	"github.com/dpb587/ssoca/auth/authz/filter"
	"github.com/dpb587/ssoca/certauth"
	envconfig "github.com/dpb587/ssoca/service/env/server/config"
)

type Config struct {
	CertAuths []CertAuthConfig `yaml:"certauths,omitempty"`
	Auth      AuthConfig       `yaml:"auth"` // deprecated in favor of env.auth_service, server.require
	Env       envconfig.Config `yaml:"env"`
	Server    ServerConfig     `yaml:"server"`
	Services  []ServiceConfig  `yaml:"services"`
}

type CertAuthConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

type ServerConfig struct {
	CertificatePath string                 `yaml:"certificate_path"`
	Host            string                 `yaml:"host"`
	Port            int                    `yaml:"port"`
	PrivateKeyPath  string                 `yaml:"private_key_path"`
	Redirect        ServerRedirectConfig   `yaml:"redirect"`
	Require         []filter.RequireConfig `yaml:"require"`
	RobotsTXT       string                 `yaml:"robotstxt"`
	TrustedProxies  ServerTrustedProxies   `yaml:"trusted_proxies"`
}

type ServerTrustedProxies []ServerTrustedProxy

func (v ServerTrustedProxies) AsIPNet() []*net.IPNet {
	var convert []*net.IPNet

	for _, r := range v {
		n := net.IPNet(r)
		convert = append(convert, &n)
	}

	return convert
}

type ServerTrustedProxy net.IPNet

var _ yaml.Unmarshaler = &ServerTrustedProxy{}

func (v *ServerTrustedProxy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	if !strings.Contains(data, "/") {
		ip := net.ParseIP(data)

		if ip.To4() != nil {
			data = fmt.Sprintf("%s/32", data)
		} else {
			data = fmt.Sprintf("%s/128", data)
		}
	}

	_, proxy, err := net.ParseCIDR(data)
	if err != nil {
		return errors.Wrap(err, "parsing trusted proxy CIDR")
	}

	*v = ServerTrustedProxy(*proxy)

	return nil
}

type ServerRedirectConfig struct {
	Root        string `yaml:"root"`
	AuthSuccess string `yaml:"auth_success"`
	AuthFailure string `yaml:"auth_failure"`
}

type AuthConfig struct {
	Type    string                 `yaml:"type"`    // deprecated
	Options map[string]interface{} `yaml:"options"` // deprecated
	Require []filter.RequireConfig `yaml:"require"` // deprecated
}

type ServiceConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Require []filter.RequireConfig `yaml:"require"`
	Options map[string]interface{} `yaml:"options"`
}

func (c *Config) ApplyDefaults() {
	c.Server.ApplyDefaults()

	for _, certauth := range c.CertAuths {
		certauth.ApplyDefaults()
	}

	for _, service := range c.Services {
		service.ApplyDefaults()
	}
}

func (c *Config) ApplyMigrations() {
	if c.Auth.Require != nil {
		if c.Server.Require != nil {
			panic("config migration failure: only one may be configured: auth.require, server.require")
		}

		c.Server.Require = c.Auth.Require
		c.Auth.Require = nil
	}

	if c.Auth.Type != "" {
		if c.Env.DefaultAuthService != "" {
			panic("config migration failure: only one may be configured: auth.type, env.default_auth_service")
		}

		c.Env.DefaultAuthService = "auth"

		c.Services = append(
			c.Services,
			ServiceConfig{
				Name:    "auth",
				Type:    fmt.Sprintf("%s_authn", c.Auth.Type),
				Options: c.Auth.Options,
			},
		)

		c.Auth.Type = ""
		c.Auth.Options = nil
	}
}

func (c *ServerConfig) ApplyDefaults() {
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}

	if c.Port == 0 {
		c.Port = 18705
	}

	if c.RobotsTXT == "" {
		c.RobotsTXT = "User-agent: *\nDisallow: /"
	}
}

func (c *CertAuthConfig) ApplyDefaults() {
	if c.Name == "" {
		c.Name = certauth.DefaultName
	}
}

func (c *ServiceConfig) ApplyDefaults() {
	if c.Name == "" {
		c.Name = c.Type
	}
}
