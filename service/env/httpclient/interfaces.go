package httpclient

import (
	"github.com/dpb587/ssoca/service/env/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetInfo() (api.InfoResponse, error)
}
