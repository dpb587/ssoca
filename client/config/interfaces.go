package config

//go:generate counterfeiter . Manager
type Manager interface {
	GetSource() string
	GetEnvironments() (EnvironmentsState, error)
	GetEnvironment(string) (EnvironmentState, error)
	SetEnvironment(EnvironmentState) error
}
