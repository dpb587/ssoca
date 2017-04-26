package management

import (
	"fmt"
	"io"
)

func SimpleCallbackHandler(_ io.Writer, data string) (ClientHandlerCallback, error) {
	fmt.Print(data)

	return nil, nil
}
