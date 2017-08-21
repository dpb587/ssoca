package profile

//go:generate counterfeiter . Manager
type Manager interface {
	Sign(data []byte) ([]byte, error)
	GetProfile() (Profile, error)
	IsCertificateValid() bool
	Renew() error
}
