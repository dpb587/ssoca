package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/ssh/httpclient"
)

type GetClient func(string) (*svchttpclient.Client, error)
