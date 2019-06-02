package installer

import (
	"errors"

	"github.com/dpb587/ssoca/service/openvpn/client/executable/finder"
	"github.com/sirupsen/logrus"
)

type MultiInstaller struct {
	Name       string
	Installers map[string]Installer
	Finder     finder.Finder
}

func (i *MultiInstaller) Install(logger logrus.FieldLogger) error {
	logger = logger.WithField("installer", i.Name)

	for method, installer := range i.Installers {
		logger.Warnf("attempting %s installation via %s", i.Name, method)

		err := installer.Install(logger)
		if err != nil {
			logger.Errorf("failed %s installation via %s: %s", i.Name, method, err)

			continue
		}

		_, _, err = i.Finder.Find()
		if err != nil {
			logger.Errorf("failed %s installation via %s: %s\n", i.Name, method, errors.New("executable still not found after installation"))

			continue
		}

		logger.Infof("completed %s installation via %s\n", i.Name, method)

		return nil
	}

	return errors.New("installation failed")
}
