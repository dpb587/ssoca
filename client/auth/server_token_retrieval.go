package auth

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/version"
)

type ServerTokenRetrieval struct {
	envURL      string
	version     version.Version
	cmdRunner   boshsys.CmdRunner
	bindAddress string
	openCommand []string
	stdout      io.Writer
	stdin       io.Reader
}

func NewServerTokenRetrieval(envURL string, ver version.Version, cmdRunner boshsys.CmdRunner, bindAddress string, openCommand []string, stdout io.Writer, stdin io.Reader) ServerTokenRetrieval {
	return ServerTokenRetrieval{
		envURL:      envURL,
		version:     ver,
		cmdRunner:   cmdRunner,
		bindAddress: bindAddress,
		openCommand: openCommand,
		stdout:      stdout,
		stdin:       stdin,
	}
}

func (str *ServerTokenRetrieval) listenForTokenCallback(tokenChannel chan string, errorChannel chan error, portChannel chan string) {
	s := &http.Server{
		Addr: str.bindAddress,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

				return
			}

			http.Redirect(w, r, fmt.Sprintf("%s%s", str.envURL, r.PostFormValue("return_to")), http.StatusTemporaryRedirect)

			tokenChannel <- r.PostFormValue("token")
		}),
	}

	err := str.listenAndServeWithPort(s, portChannel)

	if err != nil {
		errorChannel <- err
	}
}

func (str *ServerTokenRetrieval) listenAndServeWithPort(srv *http.Server, portChannel chan string) error {
	addr := srv.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(ln.Addr().String())

	portChannel <- port

	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (str *ServerTokenRetrieval) waitForTokenInput(tokenChannel chan string, errorChannel chan error) {
	for {
		fmt.Fprintf(str.stdout, "token> ")

		var token string
		_, err := fmt.Fscanf(str.stdin, "%s", &token)
		if err != nil {
			if err == io.EOF {
				// log?
				return
			}

			errorChannel <- err

			return
		}

		tokenChannel <- token

		break
	}
}

func (str *ServerTokenRetrieval) Retrieve(ctx context.Context, baseurl string) (string, error) {
	stdinChannel := make(chan string)
	tokenChannel := make(chan string)
	errorChannel := make(chan error)
	portChannel := make(chan string)

	go str.listenForTokenCallback(tokenChannel, errorChannel, portChannel)

	port := <-portChannel

	fullurl := fmt.Sprintf("%s%s?client_port=%s&client_version=%s", str.envURL, baseurl, port, url.QueryEscape(str.version.Semver))

	openCommand := str.openCommand
	foundURL := false

	for argIdx, argVal := range openCommand {
		if argVal == "((url))" {
			openCommand[argIdx] = fullurl
			foundURL = true

			break
		}
	}

	if !foundURL {
		openCommand = append(openCommand, fullurl)
	}

	openExecutable := openCommand[0]
	openCommand = openCommand[1:]

	str.cmdRunner.RunComplexCommand(boshsys.Command{
		Name: openExecutable,
		Args: openCommand,

		KeepAttached: true,
	})

	fmt.Fprintln(str.stdout, "Visit the following link to receive an authentication token...")
	fmt.Fprintln(str.stdout, "")
	fmt.Fprintln(str.stdout, "  ", fullurl)
	fmt.Fprintln(str.stdout, "")

	go str.waitForTokenInput(stdinChannel, errorChannel)

	select {
	case token := <-tokenChannel:
		fmt.Fprintln(str.stdout, "(received from browser)")
		fmt.Fprintln(str.stdout, "")

		return token, nil
	case token := <-stdinChannel:
		fmt.Fprintln(str.stdout, "")

		return token, nil
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errorChannel:
		return "", err
	}
}
