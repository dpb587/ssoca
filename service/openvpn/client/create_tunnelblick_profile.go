package client

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
)

type CreateTunnelblickProfileOpts struct {
	SsocaExec     string
	SkipAuthRetry bool

	Directory string
	FileName  string
}

func (s Service) CreateTunnelblickProfile(opts CreateTunnelblickProfileOpts) (string, error) {
	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return "", errors.Wrap(err, "getting config manager")
	}

	ssocaExec := opts.SsocaExec
	if ssocaExec == "" {
		ssocaExec = "ssoca" // TODO ssoca.exe
	}

	ssocaExec, err = exec.LookPath(ssocaExec)
	if err != nil {
		return "", errors.Wrap(err, "resolving ssoca executable")
	}

	dir := opts.Directory
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "getting working directory")
		}
	}

	file := opts.FileName
	if file == "" {
		file = s.runtime.GetEnvironmentName()

		if s.name != "openvpn" {
			file = fmt.Sprintf("%s-%s", file, s.name)
		}
	}

	dir = fmt.Sprintf("%s/%s.tblk", dir, file)

	dirAbs, err := s.fs.ExpandPath(dir)
	if err != nil {
		return "", errors.Wrap(err, "expanding path")
	}

	err = s.fs.MkdirAll(dirAbs, 0700)
	if err != nil {
		return "", errors.Wrap(err, "creating target directory")
	}

	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return "", errors.Wrap(err, "getting client")
	}

	profile, err := client.BaseProfile()
	if err != nil {
		return "", errors.Wrap(err, "getting base profile")
	}

	pathConfigOvpn := fmt.Sprintf("%s/config.ovpn", dirAbs)

	err = s.fs.WriteFileString(pathConfigOvpn, profile)
	if err != nil {
		return "", errors.Wrap(err, "writing config.ovpn")
	}

	err = s.fs.Chmod(pathConfigOvpn, 0400)
	if err != nil {
		return "", errors.Wrap(err, "chmoding config.ovpn")
	}

	pathPreConnect := fmt.Sprintf("%s/pre-connect.sh", dirAbs)

	var preconnectScriptBuf bytes.Buffer
	err = tunnelblickPreConnectScript.Execute(
		&preconnectScriptBuf,
		struct {
			Exec        string
			Config      string
			Environment string
			Service     string
		}{
			Exec:        ssocaExec,
			Config:      configManager.GetSource(),
			Environment: s.runtime.GetEnvironmentName(),
			Service:     s.name,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "generating Tunnelblick pre-connect.sh")
	}

	preconnectScript := preconnectScriptBuf.String()

	err = s.fs.WriteFileString(pathPreConnect, preconnectScript)
	if err != nil {
		return "", errors.Wrap(err, "writing pre-connect.sh")
	}

	err = s.fs.Chmod(pathPreConnect, 0500)
	if err != nil {
		return "", errors.Wrap(err, "chmoding pre-connect.sh")
	}

	pathInstall := fmt.Sprintf("%s/ssoca-install.sh", dirAbs)

	err = s.fs.WriteFileString(pathInstall, tunnelblickInstallScript)
	if err != nil {
		return "", errors.Wrap(err, "writing ssoca-install.sh")
	}

	err = s.fs.Chmod(pathInstall, 0500)
	if err != nil {
		return "", errors.Wrap(err, "chmoding ssoca-install.sh")
	}

	return dir, nil
}

var tunnelblickPreConnectScript = template.Must(template.New("script").Parse(`#!/bin/bash

set -u

REM() { /bin/echo $( date -u +"%Y-%m-%dT%H:%M:%SZ" ) "$@"; }

profile="$( basename "$( dirname "$( dirname "$( dirname "$0" )" )" )" )"
shadow="$( dirname "$0" )/config.ovpn"

file="$HOME/Library/Application Support/Tunnelblick/Configurations/$profile/Contents/Resources/config.ovpn"

REM "renewing profile"
sudo -Hnu "$USER" -- {{.Exec}} --config "{{.Config}}" --environment "{{.Environment}}" openvpn create-profile --service "{{.Service}}" > "$file.tmp"
exit=$?

if [[ "0" != "$exit" ]] || [[ ! -s "$file.tmp" ]]; then
  rm "$file.tmp"

  REM "exiting with failure"
  exit $exit
fi

set -e

REM "patching profile"
echo "remap-usr1 SIGTERM" >> "$file.tmp"

REM "installing profile"
mv -f "$file.tmp" "$file"

REM "installing shadow copy"
cat "$file" > "$shadow"

REM done
`))

var tunnelblickInstallScript = `#!/bin/bash

set -eu

[ -n "${SUDO_USER:-}" ] || ( echo "ERROR: This install script must be run with sudo" >&2 && exit 1 )

if [[ "$( ps aux | grep Tunnelblick.app | grep -v grep )" != "" ]]; then
	# Tunnelblick rewrites its preferences at exit, overwriting our Keep Connected option.
	# Profiles also do not automatically show up when adding new ones.
	echo "ERROR: Tunnelblick appears to be running. To ensure this profile is installed" >&2
	echo "       correctly, please quit the application before trying again." >&2

	exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

profile="$( basename "$DIR" )"
name="${profile%.tblk}"
preferencesFile="$HOME/Library/Preferences/net.tunnelblick.tunnelblick.plist"
profileDir="$HOME/Library/Application Support/Tunnelblick/Configurations/$profile"
shadowDir="/Library/Application Support/Tunnelblick/Users/$SUDO_USER/$profile"

#
# profile
#

# create and secure profile directory in case Tunnelblick has never been run before
if [ ! -e "$HOME/Library/Application Support/Tunnelblick/Configurations" ]; then
  mkdir -p "$HOME/Library/Application Support/Tunnelblick/Configurations"

  chown "$SUDO_USER":admin "$HOME/Library/Application Support/Tunnelblick"
  chmod 750 "$HOME/Library/Application Support/Tunnelblick"

  chown "$SUDO_USER":admin "$HOME/Library/Application Support/Tunnelblick/Configurations"
  chmod 750 "$HOME/Library/Application Support/Tunnelblick/Configurations"
fi

mkdir -p "$profileDir/Contents/Resources"

chown "$SUDO_USER":admin "$profileDir"
chmod 750 "$profileDir"

chown "$SUDO_USER":admin "$profileDir/Contents"
chmod 750 "$profileDir/Contents"

chown "$SUDO_USER":admin "$profileDir/Contents/Resources"
chmod 750 "$profileDir/Contents/Resources"

cp "$DIR/config.ovpn" "$profileDir/Contents/Resources/config.ovpn"
chown "$SUDO_USER":admin "$profileDir/Contents/Resources/config.ovpn"
chmod 740 "$profileDir/Contents/Resources/config.ovpn"

cp "$DIR/pre-connect.sh" "$profileDir/Contents/Resources/pre-connect.sh"
chown "$SUDO_USER":admin "$profileDir/Contents/Resources/pre-connect.sh"
chmod 750 "$profileDir/Contents/Resources/pre-connect.sh"

#
# shadow
#

# create and secure shadow directory in case Tunnelblick has never been run before
if [ ! -e "/Library/Application Support/Tunnelblick/Users/$SUDO_USER" ]; then
  mkdir -p "/Library/Application Support/Tunnelblick/Users/$SUDO_USER"

  chown "root:wheel" "/Library/Application Support/Tunnelblick"
  chmod 755 "/Library/Application Support/Tunnelblick"

  chown "root:wheel" "/Library/Application Support/Tunnelblick/Users"
  chmod 755 "/Library/Application Support/Tunnelblick/Users"

  chown "root:wheel" "/Library/Application Support/Tunnelblick/Users/$SUDO_USER"
  chmod 755 "/Library/Application Support/Tunnelblick/Users/$SUDO_USER"
fi

mkdir -p "$shadowDir/Contents/Resources"

chown root:wheel "$shadowDir"
chmod 755 "$shadowDir"

chown root:wheel "$shadowDir/Contents"
chmod 755 "$shadowDir/Contents"

chown root:wheel "$shadowDir/Contents/Resources"
chmod 755 "$shadowDir/Contents/Resources"

cp "$DIR/config.ovpn" "$shadowDir/Contents/Resources/config.ovpn"
chown root:wheel "$shadowDir/Contents/Resources/config.ovpn"
chmod 700 "$shadowDir/Contents/Resources/config.ovpn"

cp "$DIR/pre-connect.sh" "$shadowDir/Contents/Resources/pre-connect.sh"
chown root:wheel "$shadowDir/Contents/Resources/pre-connect.sh"
chmod 700 "$shadowDir/Contents/Resources/pre-connect.sh"

#
# preferences
#

# the configuration is generated; don't confuse people into thinking they can edit it
sudo -u "$SUDO_USER" -- defaults write net.tunnelblick.tunnelblick "$name-disableEditConfiguration" -bool yes

# since we cannot live-reload new certs here, we "unexpectedly" exit and want Tunnelblick to restart us
sudo -u "$SUDO_USER" -- defaults write net.tunnelblick.tunnelblick "$name-keepConnected" -bool yes

# by default, avoid checking if IP changed after connection to avoid alerts and monitoring
sudo -u "$SUDO_USER" -- defaults write net.tunnelblick.tunnelblick "$name-notOKToCheckThatIPAddressDidNotChangeAfterConnection" -bool yes

# be helpful and pre-select our new profile in case they want to change more options
sudo -u "$SUDO_USER" -- defaults write net.tunnelblick.tunnelblick leftNavSelectedDisplayName -string "$name"

#
# fyi
#

echo "The profile '$name' has successfully been installed."
`
