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
	"github.com/sirupsen/logrus"
)

type Config struct {
	CertAuths []CertAuthConfig `yaml:"certauths,omitempty"`
	Auth      AuthConfig       `yaml:"auth"`
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
	CertificatePath string               `yaml:"certificate_path"`
	Host            string               `yaml:"host"`
	Port            int                  `yaml:"port"`
	PrivateKeyPath  string               `yaml:"private_key_path"`
	Redirect        ServerRedirectConfig `yaml:"redirect"`
	TrustedProxies  ServerTrustedProxies `yaml:"trusted_proxies"`
	RobotsTXT       string               `yaml:"robotstxt"`
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
	AuthSuccess string `yaml:"auth_success"` // TODO deprecated
	AuthFailure string `yaml:"auth_failure"` // TODO deprecated
}

type AuthConfig struct {
	Type           string                 `yaml:"type"`    // TODO deprecated
	Options        map[string]interface{} `yaml:"options"` // TODO deprecated
	Require        []filter.RequireConfig `yaml:"require"`
	DefaultService string                 `yaml:"default_auth_service"`
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

func (c *Config) ApplyMigrations(logger logrus.FieldLogger) error {
	// introduced v0.17.0; remove v1.0.0
	if c.Auth.Type != "" {
		logger.Warn("authentication should now be configured as a service (auth.type, auth.options are deprecated)")

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

	if c.Server.Redirect.AuthFailure != "" {
		logger.Warn("authentication redirects should now be configured through service options (server.redirect.auth_failure is deprecated)")

		for srvIdx, srv := range c.Services {
			if srv.Type == "github_authn" || srv.Type == "google_authn" {
				c.Services[srvIdx].Options["redirect_failure"] = c.Server.Redirect.AuthFailure
			}
		}

		c.Server.Redirect.AuthFailure = ""
	}

	if c.Server.Redirect.AuthSuccess != "" {
		logger.Warn("authentication redirects should now be configured through service options (server.redirect.auth_failure is deprecated)")

		for srvIdx, srv := range c.Services {
			if srv.Type == "github_authn" || srv.Type == "google_authn" {
				c.Services[srvIdx].Options["redirect_success"] = c.Server.Redirect.AuthSuccess
			}
		}

		c.Server.Redirect.AuthSuccess = ""
	}

	return nil
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
