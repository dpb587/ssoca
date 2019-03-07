// +build !darwin
// +build !windows

package helper

var guessExecutableName = "openvpn"
var guessExecutablePaths = []string{
	"/usr/local/sbin/openvpn",
	"/usr/sbin/openvpn",
}
var guessExecutableSuggestions = ""
