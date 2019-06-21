package authn

type AuthorizationToken struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`

	Token string `yaml:"token,omitempty"` // deprecated; newer versions explicitly use type+value
}
