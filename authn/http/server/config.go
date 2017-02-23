package server

type Config struct {
	Users []UserConfig `yaml:"users"`
}

type UserConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Attributes map[string]interface{}
}
