package apputil

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	Number           = "0123456789"
	UpCaseAlphabet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DownCaseAlphabet = "abcdefghijklmnopqrstuvwsyz"
)

func RemoveSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func ContainsSpace(str string) bool {
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

// 伏せ字にする
func HideWords(str string) string {
	length := utf8.RuneCountInString(str)
	res := ""
	for i := 0; i < length; i++ {
		res = res + "●"
	}
	return res
}

func GenerateHash(str string) (string, error) {
	const cost = 10
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

/*
	許容文字のみを含んでいるかを判定する
	許容文字以外の文字が含まれている場合、falseとなる

	ex) ContainsAllowedStrOnly("あ", "Abc123")→ false
*/
func ContainsAllowedStrOnly(target string, allowed string) bool {
	for _, e := range target {
		if !strings.ContainsAny(string(e), allowed) {
			return false
		}
	}
	return true
}
