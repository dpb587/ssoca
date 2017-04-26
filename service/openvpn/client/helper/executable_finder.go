package helper

import (
	"errors"
	"os/exec"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

// var tunnelblickDefault = "/Applications/Tunnelblick.app/Contents/Resources/openvpn/openvpn-2.4.1-openssl-1.0.2k/openvpn"
var tunnelblickDefault = "/Applications/Tunnelblick.app/Contents/Resources/openvpn/default"

type ExecutableFinder struct {
	Paths []string
	FS    boshsys.FileSystem
}

func (ef ExecutableFinder) Find() (string, error) {
	path, err := exec.LookPath("openvpn")
	if err == nil {
		return path, nil
	}

	if ef.FS.FileExists(tunnelblickDefault) {
		return tunnelblickDefault, nil
	}

	return "", errors.New("Failed to find the openvpn executable")
}
