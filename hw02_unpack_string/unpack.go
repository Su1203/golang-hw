package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func isBab() (string, error) {
	return "", ErrInvalidString
}

func Unpack(s string) (string, error) {
	runes := []rune(s)
	var result strings.Builder

	//проверка 1 символа, что не число
	if len(runes) > 1 && unicode.IsDigit(runes[0]) {
		return isBab()
	}

	for i := 0; i < len(runes); i++ {
		char := runes[i]
		isString := false

		//экранирован
		if runes[i] == '\\' {
			if unicode.IsDigit(runes[i+1]) || runes[i+1] == '\\' {
				char = runes[i+1]
				isString = true
				i++
			} else {
				return isBab()
			}
		}
		if !unicode.IsDigit(char) || isString {
			//следующий символ - цифра
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {

				//нет - это не цифра, а число
				if i+2 < len(runes) && unicode.IsDigit(runes[i+2]) {
					fmt.Println("runes", string(char))
					return isBab()
				}

				count, err := strconv.Atoi(string(runes[i+1]))
				if err == nil {
					result.WriteString(strings.Repeat(string(char), count))
				} else {
					return isBab()
				}
			} else {
				result.WriteString(string(char))
			}
		}
	}

	return result.String(), nil
}
