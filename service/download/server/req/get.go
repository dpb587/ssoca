package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/dpb587/ssoca/server/api"
	"github.com/dpb587/ssoca/server/service/req"
	svcconfig "github.com/dpb587/ssoca/service/download/config"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Get struct {
	Paths []svcconfig.PathConfig
	FS    boshsys.FileSystem
}

var _ req.RouteHandler = Get{}

func (h Get) Route() string {
	return "get"
}

func (h Get) Execute(r *http.Request, w http.ResponseWriter) error {
	name := r.URL.Query().Get("name")
	if name == "" {
		return api.NewError(errors.New("Missing query parameter: name"), 404, "")
	}

	for _, file := range h.Paths {
		if file.Name == name {
			fh, err := h.FS.OpenFile(file.Path, os.O_RDONLY, 0)
			if err != nil {
				return bosherr.WrapError(err, "Opening file for reading")
			}

			defer fh.Close()

			w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Name))
			w.Header().Add("Content-Length", strconv.FormatInt(file.Size, 10))
			io.Copy(w, fh)

			return nil
		}
	}

	return api.NewError(fmt.Errorf("Invalid file name: %s", name), 404, "")
}
