package uaaauth

type Service struct{}

func (Service) Type() string {
	return "uaa_authn"
}

func (Service) Version() string {
	return "0.1.0"
}
