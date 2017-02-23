package auth

import "fmt"

type simpleToken struct {
	username   string
	attributes map[string]interface{}
}

func NewSimpleToken(username string, attributes map[string]interface{}) simpleToken {
	return simpleToken{
		username:   username,
		attributes: attributes,
	}
}

func (t simpleToken) Username() string {
	return t.username
}

func (t simpleToken) Attributes() map[string]interface{} {
	return t.attributes
}

func (t simpleToken) HasAttribute(name string) bool {
	_, ok := t.attributes[name]

	return ok
}

func (t simpleToken) GetAttribute(name string) (interface{}, error) {
	value, ok := t.attributes[name]
	if !ok {
		return nil, fmt.Errorf("Attribute not defined: %s", name)
	}

	return value, nil
}
