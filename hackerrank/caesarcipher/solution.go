package caesarcipher

import (
	"strings"
	"unicode"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func solution(s string, spaces int) string {
	var cipher string
	var isUpper bool

	for _, char := range s {
		if !unicode.IsLetter(char) {
			cipher += string(char)
			continue
		}

		letter := string(char)

		if unicode.IsUpper(char) {
			letter = strings.ToLower(letter)
			isUpper = true
		}

		indexInAlphabet := strings.Index(alphabet, letter)
		newPosition := indexInAlphabet + spaces

		if newPosition > (len(alphabet) - 1) {
			newPosition = newPosition - len(alphabet)
		}

		newValue := string(alphabet[newPosition])

		if isUpper {
			cipher += strings.ToUpper(newValue)
		} else {
			cipher += newValue
		}
		isUpper = false
	}
	return cipher
}
