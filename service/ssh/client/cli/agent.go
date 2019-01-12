// +build !windows

package cli

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"syscall"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	sshagent "golang.org/x/crypto/ssh/agent"

	clierrors "github.com/dpb587/ssoca/cli/errors"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/ssh/client"
)

type Agent struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	Foreground bool   `long:"foreground" description:"Stay in foreground"`
	Socket     string `long:"socket" description:"Socket path (ensure the directory has restricted permissions)"`

	serviceFactory svc.ServiceFactory
	fs             boshsys.FileSystem
	cmdRunner      boshsys.CmdRunner
}

var _ flags.Commander = Agent{}

func (c Agent) Execute(args []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	if c.Socket == "" {
		tmpdir, err := ioutil.TempDir("", "ssoca-ssh-agent")
		if err != nil {
			return errors.Wrap(err, "Creating temporary directory")
		}

		err = os.Chmod(tmpdir, 0700)
		if err != nil {
			return errors.Wrap(err, "Setting temporary directory permissions")
		}

		c.Socket = fmt.Sprintf("%s/agent.sock", tmpdir)
	} else {
		socket, err := c.fs.ExpandPath(c.Socket)
		if err != nil {
			return errors.Wrap(err, "Expanding socket path")
		}

		c.Socket = socket
	}

	if !c.Foreground && len(args) == 0 {
		configManager, err := c.Runtime.GetConfigManager()
		if err != nil {
			return errors.Wrap(err, "Getting config manager")
		}

		executable, err := os.Executable()
		if err != nil {
			return errors.Wrap(err, "Finding executable")
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
			return errors.Wrap(err, "Starting agent in background")
		}

		pid := process.Pid

		err = process.Release()
		if err != nil {
			return errors.Wrap(err, "Detaching from agent")
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
			return errors.Wrapf(err, "Connecting to current SSH agent (%s)", envAuthSock)
		}

		parentAgent = sshagent.NewClient(socket)
	}

	sshAgent := service.NewAgent(parentAgent)

	socket, err := net.Listen("unix", c.Socket)
	if err != nil {
		return errors.Wrap(err, "Opening socket")
	}

	defer socket.Close()

	pid := strconv.Itoa(os.Getpid())
	os.Setenv("SSH_AUTH_SOCK", c.Socket)
	os.Setenv("SSH_AGENT_PID", pid)

	if len(args) > 0 {
		go sshAgent.Listen(socket)

		_, _, exit, err := c.cmdRunner.RunComplexCommand(boshsys.Command{
			Name: args[0],
			Args: args[1:],

			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})

		// @todo doesn't seem to get here

		if err != nil && exit == 0 {
			return errors.Wrap(err, "Executing command")
		}

		return clierrors.Exit{Code: exit}
	}

	c.printEnv(pid)

	return sshAgent.Listen(socket)
}

func (c Agent) printEnv(pid string) {
	stdout := c.Runtime.GetStdout()

	fmt.Fprintf(stdout, "SSH_AUTH_SOCK=%s; export SSH_AUTH_SOCK;\n", c.Socket)
	fmt.Fprintf(stdout, "SSH_AGENT_PID=%s; export SSH_AGENT_PID;\n", pid)
	fmt.Fprintf(stdout, "echo ssoca ssh agent pid %s;\n", pid)
}
