package storage

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"gopkg.in/yaml.v2"
)

type YAMLFormat struct{}

var _ Storage = YAMLFormat{}

func (s YAMLFormat) Get(data string, get interface{}) error {
	err := yaml.Unmarshal([]byte(data), get)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshaling YAML")
	}

	return nil
}

func (l YAMLFormat) Put(_ string, put interface{}) (string, error) {
	out, err := yaml.Marshal(put)
	if err != nil {
		return "", bosherr.WrapError(err, "Marshaling YAML")
	}

	return string(out), nil
}
