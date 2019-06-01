package httpclient

import (
	"io"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/ssoca/service/file/api"
)

//go:generate counterfeiter . Client
type Client interface {
	GetMetadata() (api.MetadataResponse, error)
	GetList() (api.ListResponse, error)
	Download(string, io.ReadWriteSeeker, *pb.ProgressBar) error
}
