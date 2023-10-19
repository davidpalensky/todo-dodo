package util

import "regexp"

// Returns false if invalid hexcode: '#xxxxxx', where x is 0-9 or a-f
func VerifyHexcode(color string) bool {
	// Copied from StackOverflow
	regex := regexp.MustCompile("^#(?:[0-9a-fA-F]{3}){1,2}$")
	return regex.MatchString(color)
}
