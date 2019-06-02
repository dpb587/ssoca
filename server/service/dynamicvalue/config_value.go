package dynamicvalue

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
)

var configValueMissing = errors.New("no template configured")

type ConfigValue struct {
	factory Factory
	value   Value
}

var _ Value = ConfigValue{}

func NewConfigValue(factory Factory) ConfigValue {
	return ConfigValue{
		factory: factory,
	}
}

func (cv *ConfigValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	value, err := cv.factory.Create(data)
	if err != nil {
		return errors.Wrap(err, "parsing dynamic value")
	}

	cv.value = value

	return nil
}

func (cv *ConfigValue) WithDefault(value Value) {
	if cv.value != nil {
		return
	}

	cv.value = value
}

func (cv ConfigValue) Evaluate(arg0 *http.Request, arg1 *auth.Token) (string, error) {
	if cv.value == nil {
		return "", nil
	}

	return cv.value.Evaluate(arg0, arg1)
}
