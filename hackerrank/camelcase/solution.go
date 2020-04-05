package camelcase

import "unicode"

func solution(s string) (numberOfWords int) {
	for _, char := range s {
		if !unicode.IsUpper(char) {
			continue
		}
		if numberOfWords == 0 {
			numberOfWords += 2
			continue
		} else {
			numberOfWords++
		}
	}
	return
}
