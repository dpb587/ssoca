package config

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	yaml "gopkg.in/yaml.v2"
)

func RemarshalYAML(from interface{}, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return bosherr.WrapError(err, "Marshaling")
	}

	err = yaml.Unmarshal(bytes, to)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshalling")
	}

	defaultable, ok := to.(Defaultable)
	if ok {
		defaultable.ApplyDefaults()
	}

	return nil
}

func RemarshalJSON(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return bosherr.WrapError(err, "Marshaling")
	}

	err = json.Unmarshal(bytes, to)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshalling")
	}

	defaultable, ok := to.(Defaultable)
	if ok {
		defaultable.ApplyDefaults()
	}

	return nil
}
