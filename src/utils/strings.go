package utils

import (
	"bytes"
)

const strikethroughRune rune = '\u0336'

func insertInEachChar(s string, r rune) string {
	var buffer bytes.Buffer
	for _, char := range s {
		buffer.WriteRune(char)
		buffer.WriteRune(r)

	}
	return buffer.String()
}

func Strikethrough(s string) string {
	return insertInEachChar(s, strikethroughRune)
}
