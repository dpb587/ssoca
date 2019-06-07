package uaaauth

type ServiceType struct{}

func (ServiceType) Type() string {
	return "uaa_authn"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
