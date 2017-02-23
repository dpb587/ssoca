package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/openvpn/httpclient"
)

type GetClient func(string) (*svchttpclient.Client, error)
type CreateUserProfile func(string) (string, error)
