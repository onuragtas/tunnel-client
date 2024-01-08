package utils

import (
	"os"
	"strings"
)

func WriteToken(token string) {
	path, _ := os.UserHomeDir()

	os.WriteFile(path+"/.token", []byte(token), 0644)
}

func ReadToken() string {
	path, _ := os.UserHomeDir()

	token, _ := os.ReadFile(path + "/.token")
	return strings.TrimSpace(string(token))
}
