package req

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/auth"
	apierr "github.com/dpb587/ssoca/server/api/errors"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Request struct {
	ID            string
	LoggerContext logrus.Fields
	AuthToken     *auth.Token

	RawRequest  *http.Request
	RawResponse http.ResponseWriter
}

func (r *Request) ReadPayload(data interface{}) error {
	bytes, err := ioutil.ReadAll(r.RawRequest.Body)
	if err != nil {
		return bosherr.WrapError(err, "Reading request body")
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return apierr.NewError(apierr.WrapError(err, "Unmarshaling request payload"), http.StatusBadRequest, "Invalid body")
	}

	return nil
}

func (r *Request) WritePayload(data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return bosherr.WrapError(err, "Marshalling response payload")
	}

	r.RawResponse.Header().Add("Content-Type", "application/json")
	io.WriteString(r.RawResponse, string(bytes))
	io.WriteString(r.RawResponse, "\n")

	return nil
}
