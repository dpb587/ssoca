package management

import (
	"log"
	"net"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	handler      ServerHandler
	bindProtocol string
	bindAddress  string

	logger logrus.FieldLogger

	listener   net.Listener
	stopSignal chan struct{}
	stopError  chan error
}

func NewServer(handler ServerHandler, bindProtocol, bindAddress string, logger logrus.FieldLogger) Server {
	return Server{
		handler:      handler,
		bindProtocol: bindProtocol,
		bindAddress:  bindAddress,
		logger:       logger,
	}
}

func (cs *Server) Start() error {
	listener, err := net.Listen(cs.bindProtocol, cs.bindAddress)
	if err != nil {
		return bosherr.WrapError(err, "Binding")
	}

	cs.listener = listener

	go cs.listen()

	return nil
}

func (cs *Server) Stop() error {
	err := cs.listener.Close()

	return err
}

func (cs *Server) ManagementConfigValue() string {
	return strings.Join(strings.Split(cs.listener.Addr().String(), ":"), " ")
}

func (cs *Server) listen() {
	conn, err := cs.listener.Accept()

	if nil != err {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			return
		}

		log.Println(err)
	}

	logger := cs.logger.WithFields(logrus.Fields{
		"net.remote.address": conn.RemoteAddr().String(),
		"net.remote.network": conn.RemoteAddr().Network(),
	})

	logger.Info("new openvpn management connection")

	client := NewClient(conn, cs.handler, logger)
	defer client.Run()

	cs.Stop()
}
