package client

import (
	sshagent "golang.org/x/crypto/ssh/agent"
)

func (s Service) NewAgent(parent sshagent.Agent) Agent {
	return Agent{
		parent:  parent,
		service: s,
	}
}
