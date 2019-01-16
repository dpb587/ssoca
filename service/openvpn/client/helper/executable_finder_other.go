// +build !darwin
// +build !windows

package helper

var guessExecutableName = "openvpn"
var guessExecutablePaths = []string{
	"/usr/local/sbin/openvpn",
}
var guessExecutableSuggestions = ""
