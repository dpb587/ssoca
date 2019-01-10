package config

import (
	"fmt"
	"net"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth/authz/filter"
	"github.com/dpb587/ssoca/certauth"
	envconfig "github.com/dpb587/ssoca/service/env/server/config"
	yaml "gopkg.in/yaml.v2"
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
		return bosherr.WrapError(err, "Parsing trusted proxy CIDR")
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
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
	Require []filter.RequireConfig `yaml:"require"`
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

func (c *ServerConfig) ApplyDefaults() {
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}

	if c.Port == 0 {
		c.Port = 18705
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
