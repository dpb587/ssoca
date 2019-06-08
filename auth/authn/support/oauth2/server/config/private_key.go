package config

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type PrivateKey struct {
	*rsa.PrivateKey
}

var _ yaml.Unmarshaler = &PrivateKey{}

func (v *PrivateKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(data))
	if err != nil {
		return errors.Wrap(err, "parsing private key")
	}

	v.PrivateKey = privateKey

	return nil
}
