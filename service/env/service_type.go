package env

type ServiceType struct{}

func (ServiceType) Type() string {
	return "env"
}

func (ServiceType) Version() string {
	return "0.6.0"
}
