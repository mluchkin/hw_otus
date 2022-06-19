package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)
	result := ""
	if len(runes) == 0 {
		return "", nil
	}
	if !unicode.IsLetter(runes[0]) {
		return "", ErrInvalidString
	}
	for i := 0; i < len(runes); i++ {
		if i+1 == (len(runes)) {
			if unicode.IsDigit(runes[i]) {
				break
			}
			result += string(runes[i])
			break
		}
		if unicode.IsSpace(runes[i]) {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(runes[i+1]) {
			if i+2 != (len(runes)) && unicode.IsDigit(runes[i+2]) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(runes[i+1]))
			if count == 0 {
				result += strings.TrimRight(string(runes[i]), string(runes[i]))
			}
			result += strings.Repeat(string(runes[i]), count)
		} else {
			if unicode.IsDigit(runes[i]) {
				continue
			}
			result += string(runes[i])
		}
	}
	return result, nil
}
