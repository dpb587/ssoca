package httpclient

import (
	"github.com/dpb587/ssoca/service/env/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetAuth() (api.AuthResponse, error)
	GetInfo() (api.InfoResponse, error)
}
