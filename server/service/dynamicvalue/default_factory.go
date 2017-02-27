package dynamicvalue

import "strings"

type DefaultFactory struct{}

func (DefaultFactory) Create(value string) (Value, error) {
	if strings.Contains(value, "{{") {
		return CreateTemplateValue(value)
	}

	return NewStringValue(value), nil
}
