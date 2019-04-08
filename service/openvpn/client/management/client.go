package management

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	conn      net.Conn
	handler   ServerHandler
	password  string
	logger    logrus.FieldLogger
	callbacks []ServerHandlerCallback
}

func NewClient(conn net.Conn, handler ServerHandler, password string, logger logrus.FieldLogger) Client {
	return Client{
		conn:      conn,
		handler:   handler,
		password:  password,
		logger:    logger,
		callbacks: []ServerHandlerCallback{},
	}
}

func (cc *Client) Run() error {
	defer cc.conn.Close()

	reader := bufio.NewReader(cc.conn)

	if cc.password != "" {
		message, err := reader.ReadString(':')
		if err != nil {
			return err
		}

		cc.logger.Debugf("openvpn management data recv: %s", message)

		if message != "ENTER PASSWORD:" {
			return errors.New("expected first client message to request password")
		}

		cc.conn.Write([]byte(cc.password + "\n"))

		cc.callbacks = append(cc.callbacks, SuccessCallback)
	}

	for {
		cc.logger.Debug("openvpn management waiting for data")

		message, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		cc.logger.Debugf("openvpn management data recv: %s", message)

		var keepGoing bool

		if strings.HasPrefix(message, ">") {
			keepGoing = cc.handleAction(message[1:])
		} else {
			keepGoing = cc.handleReaction(message)
		}

		if !keepGoing {
			cc.logger.Info("openvpn management exiting loop")

			break
		}
	}

	return nil
}

func (cc *Client) handleAction(message string) bool {
	split := strings.SplitN(message, ":", 2)
	if len(split) != 2 {
		fmt.Println(fmt.Errorf("unexpected realtime message: %s", message))

		return true
	}

	var callback ServerHandlerCallback
	var err error

	messageSource, messageText := split[0], split[1][0:len(split[1])-1]

	switch messageSource {
	case "INFO":
		cc.logger.Info(messageText)
	case "NEED-CERTIFICATE":
		callback, err = cc.handler.NeedCertificate(cc.conn, messageText)
	case "RSA_SIGN":
		callback, err = cc.handler.SignRSA(cc.conn, messageText)
	case "FATAL":
		cc.logger.Error(messageText)

		return false
	default:
		cc.logger.Warnf("unexpected realtime message: %s", messageSource)

		return true
	}

	if err != nil {
		cc.logger.Errorf("unexpected error: %s", err)
		cc.conn.Write([]byte("signal SIGTERM\n"))
	} else if callback != nil {
		cc.callbacks = append(cc.callbacks, callback)
	}

	return true
}

func (cc *Client) handleReaction(message string) bool {
	if len(cc.callbacks) == 0 {
		cc.logger.Warnf("unexpected callback message: %s", message)

		return true
	}

	callback, err := cc.callbacks[0](cc.conn, message)
	cc.callbacks = cc.callbacks[1:]

	if err != nil {
		fmt.Println(err)
	} else if callback != nil {
		cc.callbacks = append(cc.callbacks, callback)
	}

	return true
}
