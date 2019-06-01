package installer

import "github.com/sirupsen/logrus"

type Installer interface {
	Install(logger logrus.FieldLogger) error
}
