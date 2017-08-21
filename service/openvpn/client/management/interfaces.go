package management

import "io"

//go:generate counterfeiter . ServerHandler
type ServerHandler interface {
	NeedCertificate(io.Writer, string) (ServerHandlerCallback, error)
	SignRSA(io.Writer, string) (ServerHandlerCallback, error)
}

type ServerHandlerCallback func(io.Writer, string) (ServerHandlerCallback, error)
