package certauth

import (
	"crypto/x509"
	"errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

var configValueMissing = errors.New("no certificate provider configured")

type ConfigValue struct {
	manager  Manager
	provider Provider
}

var _ Provider = ConfigValue{}

func NewConfigValue(manager Manager) ConfigValue {
	return ConfigValue{
		manager: manager,
	}
}

func (cv *ConfigValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	return cv.Configure(data)
}

func (cv *ConfigValue) Configure(name string) error {
	provider, err := cv.manager.Get(name)
	if err != nil {
		return bosherr.WrapError(err, "Getting certificate authority")
	}

	cv.provider = provider

	return nil
}

func (cv ConfigValue) IsConfigured() bool {
	return cv.provider != nil
}

func (cv ConfigValue) Name() string {
	if !cv.IsConfigured() {
		panic(configValueMissing)
	}

	return cv.provider.Name()
}

func (cv ConfigValue) GetCertificate() (*x509.Certificate, error) {
	if !cv.IsConfigured() {
		panic(configValueMissing)
	}

	return cv.provider.GetCertificate()
}

func (cv ConfigValue) GetCertificatePEM() (string, error) {
	if !cv.IsConfigured() {
		panic(configValueMissing)
	}

	return cv.provider.GetCertificatePEM()
}

func (cv ConfigValue) SignCertificate(arg0 *x509.Certificate, arg1 interface{}, arg2 logrus.Fields) ([]byte, error) {
	if !cv.IsConfigured() {
		panic(configValueMissing)
	}

	return cv.provider.SignCertificate(arg0, arg1, arg2)
}

func (cv ConfigValue) SignSSHCertificate(arg0 *ssh.Certificate, arg1 logrus.Fields) error {
	if !cv.IsConfigured() {
		panic(configValueMissing)
	}

	return cv.provider.SignSSHCertificate(arg0, arg1)
}
