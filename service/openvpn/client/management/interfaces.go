package management

import "io"

type ClientHandler interface {
	NeedCertificate(io.Writer, string) (ClientHandlerCallback, error)
	SignRSA(io.Writer, string) (ClientHandlerCallback, error)
}

type ClientHandlerCallback func(io.Writer, string) (ClientHandlerCallback, error)
