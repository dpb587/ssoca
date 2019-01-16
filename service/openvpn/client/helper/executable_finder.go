package helper

import (
	"errors"
	"os/exec"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type ExecutableFinder struct {
	Paths []string
	FS    boshsys.FileSystem
}

func (ef ExecutableFinder) Find() (string, bool, error) {
	path, err := exec.LookPath(guessExecutableName)
	if err == nil {
		return path, false, nil
	}

	for _, path := range guessExecutablePaths {
		paths, err := ef.FS.Glob(path)
		if err != nil {
			continue
		}

		for _, path := range paths {
			return path, true, nil
		}
	}

	return "", false, errors.New("Failed to find the openvpn executable")
}
