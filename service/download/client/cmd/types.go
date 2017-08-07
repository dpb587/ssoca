package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/download/httpclient"
)

type GetClient func(string, bool) (svchttpclient.Client, error)
