package config

import "github.com/dpb587/ssoca/auth"

type Config struct {
	Users []User `yaml:"users"`
}

type User struct {
	Username   string                          `yaml:"username"`
	Password   string                          `yaml:"password"`
	Groups     []string                        `yaml:"groups"`
	Attributes map[auth.TokenAttribute]*string `yaml:"attributes"`
}
