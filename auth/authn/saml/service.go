package google

type Service struct{}

func (Service) Type() string {
	return "saml_authn"
}

func (Service) Version() string {
	return "0.1.0"
}
