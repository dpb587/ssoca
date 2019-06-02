package config

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

var EnvironmentOptionAuthBind = NewStringEnvironmentOption("auth.bind", "0.0.0.0:0")
var EnvironmentOptionAuthOpenCommand = NewStringSliceEnvironmentOption("auth.open_command", defaultAuthOpenCommand)

type State struct {
	Environments EnvironmentsState `yaml:"environments,omitempty"`
}

type EnvironmentsState []EnvironmentState

type EnvironmentState struct {
	URL           string                 `yaml:"url"`
	CACertificate string                 `yaml:"ca_certificate,omitempty"`
	Alias         string                 `yaml:"alias,omitempty"`
	Auth          *EnvironmentAuthState  `yaml:"auth,omitempty"`
	Options       map[string]interface{} `yaml:"options,omitempty"`
}

func (e EnvironmentState) GetCACertificate() (*x509.Certificate, error) {
	var cert *x509.Certificate

	block, _ := pem.Decode([]byte(e.CACertificate))
	if block == nil {
		return cert, errors.New("parsing CA certificate: Missing PEM block")
	}

	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return cert, errors.New("parsing CA certificate: Not a certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return cert, errors.Wrap(err, "parsing CA certificate")
	}

	return cert, nil
}

func (e EnvironmentState) GetOption(option EnvironmentOption) error {
	val, found := e.Options[option.Key()]
	if !found {
		val = option.GetDefaultValue()
	}

	return option.SetValue(val)
}

func (e *EnvironmentState) SetOption(key string, val interface{}) {
	if e.Options == nil {
		e.Options = map[string]interface{}{}
	}

	e.Options[key] = val
}

type EnvironmentAuthState struct {
	Type    string      `yaml:"type"`
	Options interface{} `yaml:"options"`
}

func (ea EnvironmentAuthState) UnmarshalOptions(typed interface{}) error {
	bytes, err := yaml.Marshal(ea.Options)
	if err != nil {
		return errors.Wrap(err, "marshalling")
	}

	err = yaml.Unmarshal(bytes, typed)
	if err != nil {
		return errors.Wrap(err, "unmarshalling to typed options")
	}

	return nil
}
