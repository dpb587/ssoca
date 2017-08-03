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

func (ef ExecutableFinder) Find() (string, error) {
	path, err := exec.LookPath(guessExecutableName)
	if err == nil {
		return path, nil
	}

	for _, path := range guessExecutablePaths {
		if ef.FS.FileExists(path) {
			return path, nil
		}
	}

	return "", errors.New("Failed to find the openvpn executable")
}
