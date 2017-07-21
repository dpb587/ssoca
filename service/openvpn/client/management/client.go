package management

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Client struct {
	handler      ClientHandler
	bindProtocol string
	bindAddress  string

	listener    net.Listener
	activeConns *sync.WaitGroup
	stopSignal  chan struct{}
	stopError   chan error
	// not threadsafe
	callbacks []ClientHandlerCallback
}

func NewClient(handler ClientHandler, bindProtocol, bindAddress string) Client {
	return Client{
		handler:      handler,
		bindProtocol: bindProtocol,
		bindAddress:  bindAddress,
		activeConns:  &sync.WaitGroup{},
		callbacks:    []ClientHandlerCallback{},
	}
}

func (cs *Client) Start() error {
	listener, err := net.Listen(cs.bindProtocol, cs.bindAddress)
	if err != nil {
		return bosherr.WrapError(err, "Binding")
	}

	cs.listener = listener

	go cs.listen()

	return nil
}

func (cs *Client) Stop() error {
	err := cs.listener.Close()
	cs.activeConns.Wait()

	return err
}

func (cs *Client) ManagementConfigValue() string {
	return strings.Join(strings.Split(cs.listener.Addr().String(), ":"), " ")
}

func (cs *Client) listen() {
	for {
		conn, err := cs.listener.Accept()

		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Println(err)
		}

		cs.activeConns.Add(1)
		go cs.handleConnection(conn)
	}
}

func (cs *Client) handleConnection(conn net.Conn) {
	defer cs.activeConns.Done()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)

			return
		}

		if strings.HasPrefix(message, ">") {
			message = message[1:]

			split := strings.SplitN(message, ":", 2)
			if len(split) != 2 {
				fmt.Println(fmt.Errorf("unexpected realtime message: %s", message))

				continue
			}

			var callback ClientHandlerCallback
			messageSource, messageText := split[0], split[1][0:len(split[1])-1]

			switch messageSource {
			case "INFO":
				fmt.Println(messageText)
			case "NEED-CERTIFICATE":
				callback, err = cs.handler.NeedCertificate(conn, messageText)
			case "RSA_SIGN":
				callback, err = cs.handler.SignRSA(conn, messageText)
			case "FATAL":
				return
			default:
				fmt.Println(fmt.Errorf("unexpected realtime message source: %s", messageSource))

				continue
			}

			if err != nil {
				fmt.Println(fmt.Errorf("unexpected error: %s", err))
				conn.Write([]byte("signal SIGTERM\n"))

				continue
			} else if callback != nil {
				cs.callbacks = append(cs.callbacks, callback)
			}
		} else {
			if len(cs.callbacks) == 0 {
				fmt.Println(fmt.Errorf("unexpected message: %s", message))

				continue
			}

			callback, err := cs.callbacks[0](conn, message)
			cs.callbacks = cs.callbacks[1:]

			if err != nil {
				fmt.Println(err)

				continue
			} else if callback != nil {
				cs.callbacks = append(cs.callbacks, callback)
			}
		}
	}
}
