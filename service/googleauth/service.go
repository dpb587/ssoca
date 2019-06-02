package googleauth

type Service struct{}

func (Service) Type() string {
	return "google_authn"
}

func (Service) Version() string {
	return "0.1.0"
}
