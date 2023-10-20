package util

import (
	"regexp"
)

// Returns false if invalid hexcode: '#xxxxxx', where x is 0-9 or a-f
func ValidateHexcode(color string) bool {
	// Copied from StackOverflow
	regex := regexp.MustCompile("^#(?:[0-9a-fA-F]{3}){1,2}$")
	//log.Println(color)
	//log.Println(regex.MatchString(color))
	return regex.MatchString(color)
}
