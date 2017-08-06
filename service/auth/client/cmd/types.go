package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/auth/httpclient"
)

type GetClient func() (svchttpclient.Client, error)
