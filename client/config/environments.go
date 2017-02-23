package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	yaml "gopkg.in/yaml.v2"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type State struct {
	Environments EnvironmentsState `yaml:"environments,omitempty"`
}

type EnvironmentsState []EnvironmentState

type EnvironmentState struct {
	URL           string                `yaml:"url"`
	CACertificate string                `yaml:"ca_certificate,omitempty"`
	Alias         string                `yaml:"alias,omitempty"`
	Auth          *EnvironmentAuthState `yaml:"auth,omitempty"`
}

type EnvironmentAuthState struct {
	Type    string      `yaml:"type"`
	Options interface{} `yaml:"options"`
}

func (e EnvironmentState) GetCACertificate() (*x509.Certificate, error) {
	var cert *x509.Certificate

	block, _ := pem.Decode([]byte(e.CACertificate))
	if block == nil {
		return cert, errors.New("Parsing CA certificate: Missing PEM block")
	}

	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return cert, errors.New("Parsing CA certificate: Not a certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return cert, bosherr.WrapError(err, "Parsing CA certificate")
	}

	return cert, nil
}

func (ea EnvironmentAuthState) UnmarshalOptions(typed interface{}) error {
	bytes, err := yaml.Marshal(ea.Options)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling")
	}

	err = yaml.Unmarshal(bytes, typed)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshalling to typed options")
	}

	return nil
}
