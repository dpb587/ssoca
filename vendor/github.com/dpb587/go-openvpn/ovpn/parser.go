package ovpn

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func Parse(profileBytes []byte) (*Profile, error) {
	profile := Profile{}
	scanner := bufio.NewScanner(bytes.NewBuffer(profileBytes))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if len(line) == 0 {
			profile.Elements = append(profile.Elements, CommentProfileElement{})
		} else if line[0] == ';' || line[0] == '#' {
			profile.Elements = append(profile.Elements, CommentProfileElement{
				Comment: string(line[1:]),
			})
		} else if line[0] == '<' {
			embedType := strings.TrimSuffix(strings.TrimPrefix(line, "<"), ">")
			embedData := []byte{}

			closeTag := fmt.Sprintf("</%s>", embedType)

			for scanner.Scan() {
				if scanner.Text() == closeTag {
					break
				}

				embedData = append(embedData, scanner.Bytes()...)
				embedData = append(embedData, '\n')
			}

			pe, err := ParseEmbedded(embedType, embedData)
			if err != nil {
				return nil, err // TODO wrap; line
			}

			profile.Elements = append(profile.Elements, pe)
		} else {
			pe, err := ParseDirective(scanner.Bytes())
			if err != nil {
				return nil, err // TODO wrap; line
			}

			profile.Elements = append(profile.Elements, pe)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err // TODO wrap
	}

	return &profile, nil
}
