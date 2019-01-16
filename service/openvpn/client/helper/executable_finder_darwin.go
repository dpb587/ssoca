package helper

var guessExecutableName = "openvpn"
var guessExecutablePaths = []string{
	"/usr/local/sbin/openvpn",
	"/Applications/Tunnelblick.app/Contents/Resources/openvpn/openvpn-2.4.*-openssl-1.*/openvpn",
	"/Applications/Shimo.app/Contents/MacOS/openvpn",
	"/Applications/Viscosity.app/Contents/MacOS/openvpn",
}
var guessExecutableSuggestions = `
If you use Homebrew, you can install the openvpn formula...

    brew install openvpn

Alternatively, the following applications will also install openvpn...

 * Tunnelblick (https://tunnelblick.net/)
 * Shimo (https://www.shimovpn.com/)
 * Viscosity (https://www.sparklabs.com/viscosity/)
`
