package auth

import "path/filepath"

type TokenAttribute string

const (
	TokenSurnameAttribute   TokenAttribute = "surname"
	TokenGivenNameAttribute TokenAttribute = "given_name"
	TokenNameAttribute      TokenAttribute = "name"
	TokenEmailAttribute     TokenAttribute = "email"
	TokenUsernameAttribute  TokenAttribute = "username"
)

type Token struct {
	ID         string
	Groups     TokenGroups
	Attributes map[TokenAttribute]*string
}

func (t Token) Name() string {
	if t.Attributes == nil {
		return ""
	}

	return *t.Attributes[TokenNameAttribute]
}

func (t Token) Email() string {
	if t.Attributes == nil {
		return ""
	}

	return *t.Attributes[TokenEmailAttribute]
}

func (t Token) Username() string {
	if t.Attributes == nil {
		return ""
	}

	return *t.Attributes[TokenUsernameAttribute]
}

type TokenGroups []string

func (tg TokenGroups) Contains(expected string) bool {
	for _, actual := range tg {
		if expected == actual {
			return true
		}
	}

	return false
}

func (tg TokenGroups) Matches(expected string) bool {
	for _, actual := range tg {
		// note this is ignoring potential errors
		if m, _ := filepath.Match(expected, actual); m {
			return true
		}
	}

	return false
}
