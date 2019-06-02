package req

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dpb587/ssoca/auth"
	apierr "github.com/dpb587/ssoca/server/api/errors"
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
		return errors.Wrap(err, "reading request body")
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return apierr.NewError(apierr.WrapError(err, "unmarshaling request payload"), http.StatusBadRequest, "invalid body")
	}

	return nil
}

func (r *Request) WritePayload(data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshalling response payload")
	}

	r.RawResponse.Header().Add("Content-Type", "application/json")
	io.WriteString(r.RawResponse, string(bytes))
	io.WriteString(r.RawResponse, "\n")

	return nil
}
