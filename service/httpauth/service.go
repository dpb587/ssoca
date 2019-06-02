package httpauth

type Service struct{}

func (Service) Type() string {
	return "http_authn"
}

func (Service) Version() string {
	return "0.1.0"
}
