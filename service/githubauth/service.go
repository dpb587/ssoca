package githubauth

type ServiceType struct{}

func (ServiceType) Type() string {
	return "github_authn"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
