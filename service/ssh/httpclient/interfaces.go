package httpclient

import (
	"github.com/dpb587/ssoca/service/ssh/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetCAPublicKey() (api.CAPublicKeyResponse, error)
	PostSignPublicKey(api.SignPublicKeyRequest) (api.SignPublicKeyResponse, error)
}
