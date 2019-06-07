package docroot

type ServiceType struct{}

func (ServiceType) Type() string {
	return "docroot"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
