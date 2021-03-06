package httpclient

import (
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	ssocaclient "github.com/dpb587/ssoca/client"
	baseclient "github.com/dpb587/ssoca/httpclient"
)

type client struct {
	upstream baseclient.Client
	runtime  ssocaclient.Runtime
}

var _ baseclient.Client = client{}

func NewClient(upstream baseclient.Client, runtime ssocaclient.Runtime) baseclient.Client {
	return client{
		upstream: upstream,
		runtime:  runtime,
	}
}

func (c client) APIGet(url string, out interface{}) error {
	err := c.upstream.APIGet(url, out)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.APIGet(url, out)
	}

	return err
}

func (c client) APIPost(url string, in interface{}, out interface{}) error {
	err := c.upstream.APIPost(url, in, out)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.APIPost(url, in, out)
	}

	return err
}

func (c client) Get(url string) (*http.Response, error) {
	res, err := c.upstream.Get(url)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.Get(url)
	}

	return res, err
}

func (c client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	res, err := c.upstream.Post(url, contentType, body)
	if err != nil && c.attemptReauthenticate(err) == nil {
		return c.upstream.Post(url, contentType, body)
	}

	return res, err
}

func (c client) attemptReauthenticate(err error) error {
	if !strings.Contains(err.Error(), "HTTP 401") { // TODO lazy improper
		return err
	}

	configManager, err := c.runtime.GetConfigManager()
	if err != nil {
		return err
	}

	env, err := c.runtime.GetEnvironment()
	if err != nil {
		return err
	}

	var authServiceName string

	if env.Auth != nil {
		// prefer to reauthenticate with the existing auth service, if present
		authServiceName = env.Auth.Name
	}

	// this assumes runtime has a specific CLI; probably not optimal
	cmd := exec.Command(
		c.runtime.GetExec(),
		"--config", configManager.GetSource(),
		"--environment", c.runtime.GetEnvironmentName(),
		"auth",
		"--service", authServiceName,
		"login",
		"--skip-verify", // we'll be verifying in our next request
	)

	// It is weird for us to be shelling out instead of handling auth in process.
	// We fork to a new process and propagate I/O because some auth supports
	// optional stdin (e.g. server_token_retrieval.go). However, you can't really
	// do non-blocking I/O, so if no user input was provided, it would hang in
	// post-auth commands. Separate process indirectly fixes it by returning
	// ownership of stdin when its done. https://github.com/dpb587/ssoca/issues/8
	cmd.Stdin = c.runtime.GetStdin()
	// Warning: STDOUT emitted by the authentication process could corrupt the
	// output of the original command (e.g. if it's being used to generate a
	// configuration file). At the moment, the only STDOUT comes from verifying
	// the token after authenticating, so --skip-verify is avoiding the issue. It
	// might be better to redirect to STDERR, but that seems a little weird.
	cmd.Stdout = c.runtime.GetStdout()
	cmd.Stderr = c.runtime.GetStderr()

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "auth exec")
	}

	return nil
}
