package cmd

import (
	svchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
	downloadhttpclient "github.com/dpb587/ssoca/service/file/httpclient"
)

type GetClient func() (svchttpclient.Client, error)
type GetDownloadClient func(string, bool) (downloadhttpclient.Client, error)
