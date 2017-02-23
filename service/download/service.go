package download

type Service struct{}

func (Service) Type() string {
	return "download"
}

func (Service) Version() string {
	return "0.1.0"
}
