package config

import (
	"github.com/dpb587/ssoca/auth/authz/filter"
)

type Config struct {
	DefaultService string                 `yaml:"default_service"`
	Require        []filter.RequireConfig `yaml:"require"`
}
