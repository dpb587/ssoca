package nettest

import (
	"net"
	"time"
)

func NewConnection() (client, server net.Conn) {
	listener, _ := net.Listen("tcp", "127.0.0.1:0")

	go func() {
		defer listener.Close()
		server, _ = listener.Accept()
	}()

	client, _ = net.Dial("tcp", listener.Addr().String())

	time.Sleep(100 * time.Millisecond) // TODO

	return
}
