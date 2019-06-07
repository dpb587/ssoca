package auth

type ServiceType struct{}

func (ServiceType) Type() string {
	return "auth"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
