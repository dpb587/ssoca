package httpauth

type ServiceType struct{}

func (ServiceType) Type() string {
	return "http_authn"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
