package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(text string) (string, error) {
	var (
		builder    = new(strings.Builder)
		lastLetter string
		err        error
		slashRune  rune = 92
		isSlash         = false
	)
	defer builder.Reset()
	for _, char := range text {
		if unicode.IsLetter(char) {
			if lastLetter, err = unpackLetterHandler(lastLetter, builder, char); err != nil {
				return "", err
			}
		} else if unicode.IsDigit(char) {
			if lastLetter, err = unpackDigitHandler(lastLetter, char, builder); err != nil {
				return "", err
			}
		} else if char == slashRune {
			if lastLetter != "" {
				if lastLetter, err = unpackDigitHandler(lastLetter, 49, builder); err != nil {
					return "", err
				}
			}
			if !isSlash {
				isSlash = true
				continue
			}
			if isSlash {
				lastLetter += string(char)
				isSlash = false
			}
			// fmt.Printf("char: %s\n", string(char))
		}
	}
	if lastLetter != "" {
		if _, err := builder.WriteString(lastLetter); err != nil {
			return "", nil
		}
	}
	return builder.String(), nil
}

func unpackDigitHandler(lastLetter string, char rune, builder *strings.Builder) (string, error) {
	if lastLetter == "" {
		return "", ErrInvalidString
	}
	digit, err := strconv.Atoi(string(char))
	if err != nil {
		return "", err
	}
	if digit < 0 {
		return "", ErrInvalidString
	}

	subString := strings.Repeat(lastLetter, digit)
	if _, err := builder.WriteString(subString); err != nil {
		return "", err
	}
	lastLetter = ""
	return lastLetter, nil
}

func unpackLetterHandler(lastLetter string, builder *strings.Builder, char rune) (string, error) {
	if lastLetter != "" {
		if _, err := builder.WriteString(lastLetter); err != nil {
			return "", err
		}
	}
	lastLetter = string(char)
	return lastLetter, nil
}
