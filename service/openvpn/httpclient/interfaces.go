package httpclient

import (
	"github.com/dpb587/ssoca/service/openvpn/api"
)

//go:generate counterfeiter . Client
type Client interface {
	BaseProfile() (string, error)
	SignUserCSR(api.SignUserCSRRequest) (api.SignUserCSRResponse, error)
}
