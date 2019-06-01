// +build !darwin
// +build !windows

package finder

var guessExecutableName = "openvpn"
var guessExecutablePaths = []string{
	"/usr/local/sbin/openvpn",
	"/usr/sbin/openvpn",
}
