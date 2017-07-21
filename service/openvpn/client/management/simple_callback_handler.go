package management

import (
	"fmt"
	"io"
	"strings"
)

func SimpleCallbackHandler(_ io.Writer, data string) (ClientHandlerCallback, error) {
	split := strings.SplitN(data, ": ", 2)
	if split[0] != "SUCCESS" {
		return nil, fmt.Errorf("Bad management command result: %s", data)
	}

	return nil, nil
}
