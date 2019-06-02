package githubauth

type Service struct{}

func (Service) Type() string {
	return "github_authn"
}

func (Service) Version() string {
	return "0.1.0"
}
