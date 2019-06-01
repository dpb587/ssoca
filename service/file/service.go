package file

type Service struct{}

func (Service) Type() string {
	return "file"
}

func (Service) Version() string {
	return "0.1.0"
}
