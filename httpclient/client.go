package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func NewClient(endpoint string, tlsClientConfig *tls.Config) *Client {
	return &Client{
		endpoint: endpoint,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:     tlsClientConfig,
				Proxy:               http.ProxyFromEnvironment,
				TLSHandshakeTimeout: 30 * time.Second,
				DisableKeepAlives:   true,
			},
		},
	}
}

type Client struct {
	*http.Client

	endpoint string
}

func (c Client) ExpandURI(uri string) string {
	if !strings.HasPrefix(uri, "/") {
		return uri
	}

	return fmt.Sprintf("%s%s", c.endpoint, uri)
}

func (c Client) APIGet(url string, out interface{}) error {
	response, err := c.Get(c.ExpandURI(url))
	if err != nil {
		return bosherr.WrapError(err, "Executing request")
	}

	return c.apiReadResponse(response, out)
}

func (c Client) APIPost(url string, out interface{}, in interface{}) error {
	requestBody, err := json.Marshal(in)
	if err != nil {
		return bosherr.WrapError(err, "Marshaling request body")
	}

	response, err := c.Post(
		c.ExpandURI(url),
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return bosherr.WrapError(err, "Executing request")
	}

	return c.apiReadResponse(response, out)
}

func (c Client) apiReadResponse(res *http.Response, out interface{}) error {
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
