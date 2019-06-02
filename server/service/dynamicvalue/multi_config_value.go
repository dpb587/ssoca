package dynamicvalue

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
)

type MultiConfigValue struct {
	factory Factory
	values  []Value
}

var _ MultiValue = MultiConfigValue{}

func NewMultiConfigValue(factory Factory) MultiConfigValue {
	return MultiConfigValue{
		factory: factory,
	}
}

func (mcv *MultiConfigValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dataSlice []string
	if err := unmarshal(&dataSlice); err != nil {
		return err
	}

	for _, data := range dataSlice {
		value, err := mcv.factory.Create(data)
		if err != nil {
			return errors.Wrap(err, "parsing dynamic value")
		}

		mcv.values = append(mcv.values, value)
	}

	return nil
}

func (mcv MultiConfigValue) Evaluate(arg0 *http.Request, arg1 *auth.Token) ([]string, error) {
	values := []string{}

	for _, value := range mcv.values {
		res, err := value.Evaluate(arg0, arg1)
		if err != nil {
			return nil, errors.Wrap(err, "evaluating template")
		}

		values = append(values, res)
	}

	return values, nil
}
