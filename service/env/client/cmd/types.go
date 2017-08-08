package cmd

import (
	downloadhttpclient "github.com/dpb587/ssoca/service/download/httpclient"
	svchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
)

type GetClient func() (svchttpclient.Client, error)
type GetDownloadClient func(string, bool) (downloadhttpclient.Client, error)
