package auth

//go:generate counterfeiter . Token
type Token interface {
	Username() string

	Attributes() map[string]interface{}
	HasAttribute(string) bool
	GetAttribute(string) (interface{}, error)
}
