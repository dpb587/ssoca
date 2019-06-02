package certauth

import "fmt"

type defaultManager struct {
	certAuths map[string]Provider
}

var _ Manager = defaultManager{}

func NewDefaultManager() Manager {
	res := defaultManager{}
	res.certAuths = map[string]Provider{}

	return res
}

func (f defaultManager) Add(certAuth Provider) {
	f.certAuths[certAuth.Name()] = certAuth
}

func (f defaultManager) Get(name string) (Provider, error) {
	provider, ok := f.certAuths[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized name: %s", name)
	}

	return provider, nil
}
