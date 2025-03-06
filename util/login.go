package util

import (
	"regexp"
)

var (
	phoneRegex = regexp.MustCompile(`^(\+?86)?1[3-9]\d{9}$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-z]{2,}$`)
)

func IsPhone(identifier string) bool {
	return phoneRegex.MatchString(identifier)
}

func IsEmaill(identifier string) bool {
	return emailRegex.MatchString(identifier)
}
