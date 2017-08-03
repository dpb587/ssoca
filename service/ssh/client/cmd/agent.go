// +build !windows

package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"syscall"

	"github.com/dpb587/ssoca/cli/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh/agent"
	"github.com/jessevdk/go-flags"
	sshagent "golang.org/x/crypto/ssh/agent"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Agent struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	Foreground bool   `long:"foreground" description:"Stay in foreground"`
	Socket     string `long:"socket" description:"Socket path (ensure the directory has restricted permissions)"`

	GetClient GetClient
	CmdRunner boshsys.CmdRunner
	FS        boshsys.FileSystem
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
	} else {
		socket, err := c.FS.ExpandPath(c.Socket)
		if err != nil {
			return bosherr.WrapError(err, "Expanding socket path")
		}

		c.Socket = socket
	}

	// do this early to detect misconfiguration before detachinig
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	if !c.Foreground && len(args) == 0 {
		configManager, err := c.Runtime.GetConfigManager()
		if err != nil {
			return bosherr.WrapError(err, "Getting config manager")
		}

		executable, err := os.Executable()
		if err != nil {
			return bosherr.WrapError(err, "Finding executable")
		}

		process, err := os.StartProcess(
			executable,
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

		c.printEnv(strconv.Itoa(pid))

		return nil
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

	defer socket.Close()

	pid := strconv.Itoa(os.Getpid())
	os.Setenv("SSH_AUTH_SOCK", c.Socket)
	os.Setenv("SSH_AGENT_PID", pid)

	if len(args) > 0 {
		go c.serveAgent(sshAgent, socket)

		_, _, exit, err := c.CmdRunner.RunComplexCommand(boshsys.Command{
			Name: args[0],
			Args: args[1:],

			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})

		// @todo doesn't seem to get here

		if err != nil && exit == 0 {
			return bosherr.WrapError(err, "Executing command")
		}

		return errors.Exit{Code: exit}
	}

	c.printEnv(pid)

	c.serveAgent(sshAgent, socket)

	return nil
}

func (c Agent) serveAgent(agent sshagent.Agent, socket net.Listener) {
	for {
		handle, err := socket.Accept()
		if err != nil {
			return
		}

		go func(handle net.Conn) {
			defer handle.Close()

			_ = sshagent.ServeAgent(agent, handle)
		}(handle)
	}
}

func (c Agent) printEnv(pid string) {
	stdout := c.Runtime.GetStdout()

	fmt.Fprintf(stdout, "SSH_AUTH_SOCK=%s; export SSH_AUTH_SOCK;\n", c.Socket)
	fmt.Fprintf(stdout, "SSH_AGENT_PID=%s; export SSH_AGENT_PID;\n", pid)
	fmt.Fprintf(stdout, "echo ssoca ssh agent pid %s;\n", pid)
}
