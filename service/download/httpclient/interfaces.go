package httpclient

import (
	"io"

	"github.com/dpb587/ssoca/service/download/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetList() (api.ListResponse, error)
	Download(string, io.ReadWriteSeeker) error
}
