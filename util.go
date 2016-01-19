//
// util.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

func skipChar(input []byte, start int, ch byte) int {
	i := start
	length := len(input)
	for i < length && input[i] == ch {
		i++
	}
	return i
}

func skipUntilChar(input []byte, start int, ch byte) int {
	i := start
	length := len(input)
	for i < length && input[i] != ch {
		i++
	}
	return i
}

func escapeSingleChar(ch byte) (string, bool) {
	if ch == '"' {
		return "&quot;", true
	}
	if ch == '&' {
		return "&amp;", true
	}
	if ch == '<' {
		return "&lt;", true
	}
	if ch == '>' {
		return "&gt;", true
	}

	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}

	if org < len(src) {
		out.Write(src[org:])
	}
}

// expandTabs replace tab characters with spaces. aligning to the next TAB_SIZE column.
// always ends output with a newline
func expandTabs(out *bytes.Buffer, line []byte, tabSize int) {
	// first, check fro common cases: no tabs, or noly tabs at beginning of line
	i, prefix := 0, 0
	slowcase := false
	for i := 0; i < len(line); i++ {
		if line[i] == '\t' {
			if prefix == i {
				prefix++
			} else {
				slowcase = true
				break
			}
		}
	}

	// no need to decode runes if all tabs are at the beginning of the line
	if !slowcase {
		for i = 0; i < prefix*tabSize; i++ {
			out.WriteByte(' ')
		}
		out.Write(line[prefix:])
		return
	}

	// the slow case: we need to count runes to figure out how
	// many spaces to insert for each tab
	column := 0
	i = 0

	for i < len(line) {
		start := i
		for i < len(line) && line[i] != '\t' {
			_, size := utf8.DecodeRune(line[i:])
			i += size
			column++
		}

		if i > start {
			out.Write(line[start:i])
		}

		if i >= len(line) {
			break
		}

		for {
			out.WriteByte(' ')
			column++
			if column%tabSize == 0 {
				break
			}
		}

		i++
	}
}

// check if the specified position is preceded by an odd number of backslashes
func isBackslashEscaped(data []byte, i int) bool {
	backslashes := 0
	for i-backslashes >= 0 && data[i-backslashes-1] == '\\' {
		backslashes++
	}

	return backslashes&1 == 1
}

// SanitizedString returns a sanitized string for the given text.
func SanitizedString(text string) string {
	var anchorName []rune
	var futureDash = false

	for _, ch := range []rune(text) {
		switch {
		case unicode.IsLetter(ch) || unicode.IsNumber(ch):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(ch))
		default:
			futureDash = true
		}
	}

	return string(anchorName)
}

// isspace test if a character is a whitespace character
func isspace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// isletter test if a character is a letter
func isletter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isalnum test if a character is a letter or a digit
func isalnum(c byte) bool {
	return (c >= '0' && c <= '9') || isletter(c)
}

// ispunct test if a character is a puncuation symbol
func ispunct(c byte) bool {
	for _, r := range []byte("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		if c == r {
			return true
		}
	}

	return false
}

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}
