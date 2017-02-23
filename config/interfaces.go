package config

//go:generate counterfeiter . Defaultable
type Defaultable interface {
	ApplyDefaults()
}
