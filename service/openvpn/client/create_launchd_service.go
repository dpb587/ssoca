package client

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type CreateLaunchdServiceOpts struct {
	SsocaExec string

	Directory string
	Name      string

	OpenvpnExec string
	RunAtLoad   bool
	LogDir      string
}

func (s Service) CreateLaunchdService(opts CreateLaunchdServiceOpts) (string, string, error) {
	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return "", "", bosherr.WrapError(err, "Getting config manager")
	}

	ssocaExec := opts.SsocaExec
	if ssocaExec == "" {
		ssocaExec = "ssoca" // TODO ssoca.exe
	}

	ssocaExec, err = exec.LookPath(ssocaExec)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Resolving ssoca executable")
	}

	openvpnExec := opts.OpenvpnExec
	if openvpnExec == "" {
		openvpnExec = "openvpn" // TODO openvpn.exe
	}

	openvpnExec, err = exec.LookPath(openvpnExec)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Resolving openvpn executable")
	}

	dir := opts.Directory
	if dir == "" {
		dir = "~/Library/LaunchAgents"
	}

	dirAbs, err := s.fs.ExpandPath(dir)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Expanding path")
	}

	logDir := opts.LogDir
	if logDir == "" {
		logDir = "~/Library/Logs"
	}

	logDir, err = s.fs.ExpandPath(logDir)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Expanding log directory")
	}

	name := opts.Name
	if name == "" {
		name = s.runtime.GetEnvironmentName()

		if s.name != "openvpn" {
			name = fmt.Sprintf("%s.%s", s.name, name)
		}

		name = fmt.Sprintf("%s.openvpn.ssoca.dpb587.github.io", name)
	}

	err = s.fs.MkdirAll(dirAbs, 0700)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Creating target directory")
	}

	pathLaunchdService := path.Join(dirAbs, fmt.Sprintf("%s.plist", name))

	var launchdServiceBuf bytes.Buffer
	err = launchdService.Execute(
		&launchdServiceBuf,
		struct {
			Name        string
			Exec        string
			Config      string
			Environment string
			Service     string
			OpenvpnExec string
			RunAtLoad   string
			LogDir      string
		}{
			Name:        name,
			Exec:        ssocaExec,
			Config:      configManager.GetSource(),
			Environment: s.runtime.GetEnvironmentName(),
			Service:     s.name,
			OpenvpnExec: openvpnExec,
			RunAtLoad:   strconv.FormatBool(opts.RunAtLoad),
			LogDir:      logDir,
		},
	)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Generating launchd service")
	}

	launchdService := launchdServiceBuf.String()

	err = s.fs.WriteFileString(pathLaunchdService, launchdService)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Writing service plist")
	}

	err = s.fs.Chmod(pathLaunchdService, 0744)
	if err != nil {
		return "", "", bosherr.WrapError(err, "Chmoding service plist")
	}

	return pathLaunchdService, name, nil
}

var launchdService = template.Must(template.New("script").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>{{.Name}}</string>
        <key>ProgramArguments</key>
        <array>
            <string>{{.Exec}}</string>
            <string>--config={{.Config}}</string>
            <string>--environment={{.Environment}}</string>
            <string>openvpn</string>
            <string>exec</string>
            <string>--service={{.Service}}</string>
            <string>--exec={{.OpenvpnExec}}</string>
        </array>
        <key>StandardErrorPath</key>
        <string>{{.LogDir}}/{{.Name}}.stderr.log</string>
        <key>StandardOutPath</key>
        <string>{{.LogDir}}/{{.Name}}.stdout.log</string>
        <key>RunAtLoad</key>
        <{{.RunAtLoad}}/>
        <key>OnDemand</key>
        <false/>
        <key>KeepAlive</key>
        <true/>
    </dict>
</plist>
`))
