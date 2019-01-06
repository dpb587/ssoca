package auth

import (
	"bytes"
	"fmt"
	"html/template"
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

var htmlReceiverPostback = template.Must(template.New("html").Parse(`
	<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>ssoca</title>
			<script>
				window.addEventListener(
					"message",
					function (event)
					{
						if (event.origin != "{{.EnvURL}}") {
					    return;
						} else if (event.data.action != "set_token") {
							return;
						}

						var http = new XMLHttpRequest();
						http.open('POST', "/", true);
						http.setRequestHeader("content-type", "application/x-www-form-urlencoded");
						http.onreadystatechange = function() {
							if (http.readyState != 4) {
								return;
							} else if (http.status < 200 || http.status >= 300) {
								return;
							}

							event.source.postMessage({"action":"done"}, "{{.EnvURL}}");
						};
						http.send("token=" + encodeURIComponent(event.data.params.token));
					},
					false
				);
			</script>
		</head>
		<body>
			Waiting for token&hellip;
		</body>
	</html>
`))

func (str *ServerTokenRetrieval) listenForTokenCallback(tokenChannel chan string, errorChannel chan error, portChannel chan string) {
	s := &http.Server{
		Addr: str.bindAddress,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token, returnTo string

			switch r.Method {
			case "GET":
				var buf bytes.Buffer

				err := htmlReceiverPostback.Execute(
					&buf,
					struct {
						EnvURL string
					}{
						EnvURL: str.envURL,
					},
				)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)

					return
				}

				w.Header().Add("content-type", "text/html")
				w.WriteHeader(200)
				w.Write(buf.Bytes())

				return
			case "POST":
				token = r.PostFormValue("token")
				returnTo = r.PostFormValue("return_to")

				if token == "" {
					http.Error(w, "Missing token", http.StatusBadRequest)
				}

				tokenChannel <- token

				if returnTo != "" {
					http.Redirect(w, r, fmt.Sprintf("%s%s", str.envURL, returnTo), http.StatusTemporaryRedirect)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
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

func (str *ServerTokenRetrieval) Retrieve(baseurl string) (string, error) {
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
	case err := <-errorChannel:
		return "", err
	}
}
