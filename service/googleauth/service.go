package googleauth

type ServiceType struct{}

func (ServiceType) Type() string {
	return "google_authn"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
