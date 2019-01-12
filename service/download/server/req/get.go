package server

import (
	"fmt"
	"io"
	"os"
	"strconv"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"

	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/service/req"
	svcconfig "github.com/dpb587/ssoca/service/download/server/config"
)

type Get struct {
	Paths []svcconfig.PathConfig
	FS    boshsys.FileSystem

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Get{}

func (h Get) Route() string {
	return "get"
}

func (h Get) Execute(request req.Request) error {
	name := request.RawRequest.URL.Query().Get("name")
	if name == "" {
		return apierr.NewError(errors.New("Missing query parameter: name"), 404, "")
	}

	for _, file := range h.Paths {
		if file.Name == name {
			fh, err := h.FS.OpenFile(file.Path, os.O_RDONLY, 0)
			if err != nil {
				return errors.Wrap(err, "Opening file for reading")
			}

			defer fh.Close()

			request.RawResponse.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Name))
			request.RawResponse.Header().Add("Content-Length", strconv.FormatInt(file.Size, 10))
			io.Copy(request.RawResponse, fh)

			return nil
		}
	}

	return apierr.NewError(fmt.Errorf("Invalid file name: %s", name), 404, "")
}
