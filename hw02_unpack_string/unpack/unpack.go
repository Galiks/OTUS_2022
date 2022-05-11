package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack распаковывает строку
// 	Пример: "a4bc2d5e" -> "aaaabccddddde"
func Unpack(text string) (string, error) {
	var (
		builder    = new(strings.Builder)
		lastLetter string
		err        error
		slashRune  rune = 92
		isSlash         = false
	)
	if len(text) == 0 {
		return text, nil
	}
	if text[len(text)-1] == byte(slashRune) {
		return "", ErrInvalidString
	}
	defer builder.Reset()
	for _, char := range text {
		if isSlash {
			lastLetter = string(char)
			isSlash = false
			continue
		}
		switch {
		case unicode.IsDigit(char):
			if err = unpackDigitHandler(lastLetter, char, builder); err != nil {
				return "", err
			}
			lastLetter = ""
		case char == slashRune:
			{
				if lastLetter != "" {
					var oneRune rune = 49
					if err = unpackDigitHandler(lastLetter, oneRune, builder); err != nil {
						return "", err
					}
					lastLetter = ""
				}
				if !isSlash {
					isSlash = true
				} else if isSlash {
					lastLetter += string(char)
					isSlash = false
				}
			}
		default:
			if lastLetter, err = unpackLetterHandler(lastLetter, char, builder); err != nil {
				return "", err
			}
		}
	}
	if lastLetter != "" {
		if _, err := builder.WriteString(lastLetter); err != nil {
			return "", err
		}
	}
	return builder.String(), nil
}

func unpackDigitHandler(lastLetter string, char rune, builder *strings.Builder) error {
	if lastLetter == "" {
		return ErrInvalidString
	}
	digit, err := strconv.Atoi(string(char))
	if err != nil {
		return err
	}
	if digit < 0 {
		return ErrInvalidString
	}

	subString := strings.Repeat(lastLetter, digit)
	if _, err := builder.WriteString(subString); err != nil {
		return err
	}
	return nil
}

func unpackLetterHandler(lastLetter string, char rune, builder *strings.Builder) (string, error) {
	if lastLetter != "" {
		if _, err := builder.WriteString(lastLetter); err != nil {
			return "", err
		}
	}
	return string(char), nil
}
