package management

import (
	"fmt"
	"io"
	"strings"
)

func SuccessCallback(_ io.Writer, data string) (ServerHandlerCallback, error) {
	split := strings.SplitN(data, ": ", 2)
	if split[0] != "SUCCESS" {
		return nil, fmt.Errorf("Bad management command result: %s", data)
	}

	return nil, nil
}
