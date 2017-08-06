package httpclient

import (
	"github.com/dpb587/ssoca/service/auth/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetInfo() (api.InfoResponse, error)
}
