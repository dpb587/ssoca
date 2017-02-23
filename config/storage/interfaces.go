package storage

//go:generate counterfeiter . Storage
type Storage interface {
	Get(string, interface{}) error
	Put(string, interface{}) (string, error)
}
