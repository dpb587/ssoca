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
)

type client struct {
	client   *http.Client
	endpoint string
}

var _ Client = client{}

func NewClient(endpoint string, goclient *http.Client) Client {
	return &client{
		endpoint: endpoint,
		client:   goclient,
	}
}

func (c client) Get(url string) (*http.Response, error) {
	res, err := c.client.Get(c.expandURI(url))
	if err == nil && res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d", res.StatusCode)
	}

	return res, err
}

func (c client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	res, err := c.client.Post(c.expandURI(url), contentType, body)
	if err == nil && res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d", res.StatusCode)
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
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, body)
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
