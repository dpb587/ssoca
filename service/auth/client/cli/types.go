package cli

import (
	envsvchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
)

type GetClient func() (envsvchttpclient.Client, error)
