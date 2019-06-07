package file

type ServiceType struct{}

func (ServiceType) Type() string {
	return "file"
}

func (ServiceType) Version() string {
	return "0.1.0"
}
