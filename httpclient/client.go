package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/version"
)

var baseUserAgent = "httpclient.ssoca.dpb587.github.com/1.0"

type client struct {
	client   *http.Client
	endpoint string
	version  version.Version
}

var _ Client = client{}

func NewClient(goclient *http.Client, v version.Version, endpoint string) Client {
	return &client{
		endpoint: endpoint,
		client:   goclient,
		version:  v,
	}
}

func (c client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", fmt.Sprintf("%s %s", c.version.Version(), baseUserAgent))

	return c.client.Do(req)
}

func (c client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.expandURI(url), nil)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating request")
	}

	res, err := c.do(req)
	if err == nil && res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %s", res.Status)
	}

	return res, err
}

func (c client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.expandURI(url), body)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating request")
	}

	req.Header.Set("Content-Type", contentType)

	res, err := c.do(req)
	if err == nil && res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %s", res.Status)
	}

	return res, err
}

func (c client) APIGet(url string, out interface{}) error {
	response, err := c.Get(url)
	if err != nil {
		return bosherr.WrapError(err, "Executing request")
	}

	return c.apiReadResponse(response, out)
}

func (c client) APIPost(url string, out interface{}, in interface{}) error {
	requestBody, err := json.Marshal(in)
	if err != nil {
		return bosherr.WrapError(err, "Marshaling request body")
	}

	response, err := c.Post(
		c.expandURI(url),
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return bosherr.WrapError(err, "Executing request")
	}

	return c.apiReadResponse(response, out)
}

func (c client) apiReadResponse(res *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return bosherr.WrapError(err, "Reading response body")
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("HTTP %s: %s", res.Status, body)
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshaling response body")
	}

	return nil
}

func (c client) expandURI(uri string) string {
	if !strings.HasPrefix(uri, "/") {
		return uri
	}

	return fmt.Sprintf("%s%s", c.endpoint, uri)
}
