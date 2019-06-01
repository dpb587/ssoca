package installer

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type MessageInstaller struct {
	Output  io.Writer
	Message string
}

func (i *MessageInstaller) Install(_ logrus.FieldLogger) error {
	fmt.Fprintf(
		i.Output,
		"%s\n  %s\n%s\n",
		strings.Repeat("=", 80),
		strings.Replace(i.Message, "\n", "\n  ", -1),
		strings.Repeat("=", 80),
	)

	return nil
}
