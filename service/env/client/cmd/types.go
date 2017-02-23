package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
)

type GetClient func() (*svchttpclient.Client, error)
