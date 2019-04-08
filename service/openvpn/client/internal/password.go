package internal

import "math/rand"

const passwordCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GeneratePassword(size int) string {
	password := make([]byte, size)

	for i := range password {
		password[i] = passwordCharacters[rand.Intn(len(passwordCharacters))]
	}

	return string(password)
}
