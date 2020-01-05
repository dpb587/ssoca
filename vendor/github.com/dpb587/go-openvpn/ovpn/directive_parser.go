package ovpn

import (
	"regexp"
)

var naiveWhitespaceRE = regexp.MustCompile(`\s+`)

func ParseDirective(raw []byte) (ProfileElement, error) {
	split := naiveWhitespaceRE.Split(string(raw), -1)

	pe := GenericDirectiveProfileElement{directive: split[0]}

	if len(split) > 1 {
		pe.args = split[1:]
	}

	return pe, nil
}
