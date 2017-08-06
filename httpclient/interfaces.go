package httpclient

import (
	"io"
	"net/http"
)

//go:generate counterfeiter . Client
type Client interface {
	APIGet(string, interface{}) error
	APIPost(string, interface{}, interface{}) error

	Get(string) (*http.Response, error)
	Post(string, string, io.Reader) (*http.Response, error)
}
