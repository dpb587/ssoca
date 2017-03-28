package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"syscall"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh/agent"
	"github.com/jessevdk/go-flags"
	sshagent "golang.org/x/crypto/ssh/agent"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Agent struct {
	clientcmd.ServiceCommand

	Foreground bool   `long:"foreground" description:"Stay in foreground"`
	Socket     string `long:"socket" description:"Socket path (ensure the directory has restricted permissions)"`

	GetClient GetClient
}

var _ flags.Commander = Agent{}

func (c Agent) Execute(args []string) error {
	if c.Socket == "" {
		tmpdir, err := ioutil.TempDir("", "ssoca-ssh-agent")
		if err != nil {
			return bosherr.WrapError(err, "Creating temporary directory")
		}

		err = os.Chmod(tmpdir, 0700)
		if err != nil {
			return bosherr.WrapError(err, "Setting temporary directory permissions")
		}

		c.Socket = fmt.Sprintf("%s/agent.sock", tmpdir)
	}

	if !c.Foreground {
		configManager, err := c.Runtime.GetConfigManager()
		if err != nil {
			return bosherr.WrapError(err, "Getting config manager")
		}

		process, err := os.StartProcess(
			os.Args[0],
			[]string{
				fmt.Sprintf("--config=%s", configManager.GetSource()),
				fmt.Sprintf("--environment=%s", c.Runtime.GetEnvironmentName()),
				"ssh",
				"agent",
				fmt.Sprintf("--service=%s", c.ServiceName),
				fmt.Sprintf("--socket=%s", c.Socket),
				"--foreground",
			},
			&os.ProcAttr{
				Env:   os.Environ(),
				Files: []*os.File{nil, nil, nil},
				Sys: &syscall.SysProcAttr{
					Setpgid: true,
				},
			},
		)
		if err != nil {
			return bosherr.WrapError(err, "Starting agent in background")
		}

		pid := process.Pid

		err = process.Release()
		if err != nil {
			return bosherr.WrapError(err, "Detaching from agent")
		}

		c.printEnv(pid)

		return nil
	}

	client, err := c.GetClient(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	var parentAgent sshagent.Agent
	envAuthSock := os.Getenv("SSH_AUTH_SOCK")

	if envAuthSock == "" {
		parentAgent = sshagent.NewKeyring()
	} else {
		socket, err := net.Dial("unix", envAuthSock)
		if err != nil {
			return bosherr.WrapErrorf(err, "Connecting to current SSH agent", envAuthSock)
		}

		parentAgent = sshagent.NewClient(socket)
	}

	sshAgent := agent.NewAgent(parentAgent, client)

	socket, err := net.Listen("unix", c.Socket)
	if err != nil {
		return bosherr.WrapError(err, "Opening socket")
	}

	c.printEnv(os.Getpid())

	for {
		handle, _ := socket.Accept()
		go func(handle net.Conn) {
			defer handle.Close()

			_ = sshagent.ServeAgent(sshAgent, handle)
		}(handle)
	}

	return nil
}

func (c Agent) printEnv(pid int) {
	stdout := c.Runtime.GetStdout()

	fmt.Fprintf(stdout, "SSH_AUTH_SOCK=%s; export SSH_AUTH_SOCK;\n", c.Socket)
	fmt.Fprintf(stdout, "SSH_AGENT_PID=%d; export SSH_AGENT_PID;\n", pid)
	fmt.Fprintf(stdout, "echo ssoca agent pid %d;\n", pid)
}
