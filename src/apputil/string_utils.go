package apputil

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

func RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func ContainsSpaces(str string) bool {
	for _, r := range str {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// https://cipepser.hatenablog.com/entry/2017/07/29/083729
func QueryEncoding(str string) string {
	str = url.QueryEscape(str)
	str = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(str, "$1%20")
	return str
}
