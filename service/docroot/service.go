package docroot

type Service struct{}

func (Service) Type() string {
	return "docroot"
}

func (Service) Version() string {
	return "0.1.0"
}
