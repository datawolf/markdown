//
// block.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"bytes"
	//"fmt"
)

// Parse block-level data.
// Note: this function and many that it calls assume that
// the input buffer ends with a newline
func (p *parser) block(out *bytes.Buffer, input []byte) {
	if len(input) == 0 || input[len(input)-1] != '\n' {
		panic("block input is missing terminating newline")
	}

	// this is called recursively: enforce a maximum depth
	if p.nesting >= p.maxNesting {
		return
	}
	p.nesting++

	// parse out one block-level construct at a time
	for len(input) > 0 {
		// prefixed header
		//
		// # Header 1
		// ## Header 2
		// ...
		// ###### Header 6
		if p.isPrefixHeader(input) {
			input = input[p.prefixHeader(out, input):]
			continue
		}

		// blank lines. note: returns the # of bytes to skip
		if i := p.isEmpty(input); i > 0 {
			input = input[i:]
			continue
		}
	}

	p.nesting--
}

func (p *parser) isPrefixHeader(input []byte) bool {
	if input[0] != '#' {
		return false
	}

	if p.flags&EXTENSION_SPACE_HEADERS != 0 {
		level := 0
		for level < 6 && level < len(input) && input[level] == '#' {
			level++
		}
		if input[level] != ' ' {
			return false
		}
	}
	return true
}

func (p *parser) prefixHeader(out *bytes.Buffer, input []byte) int {
	level := 0
	for level < 6 && input[level] == '#' {
		level++
	}

	start := skipChar(input, level, ' ')
	end := skipUntilChar(input, start, '\n')
	skip := end
	id := ""

	for end > 0 && input[end-1] == '#' {
		if isBackslashEscaped(input, end-1) {
			break
		}
		end--
	}
	for end > 0 && input[end-1] == ' ' {
		end--
	}

	if end > start {
		if id == "" && p.flags&EXTENSION_AUTO_HEADER_IDS != 0 {
			id = SanitizedString(string(input[start:end]))
		}
		work := func() bool {
			p.inline(out, input[start:end])
			return true
		}
		p.r.Header(out, work, level, id)
	}

	return skip
}

func (p *parser) isEmpty(data []byte) int {
	// it is okay to call isEmpty on an empty buffer
	if len(data) == 0 {
		return 0
	}

	var i int
	for i = 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] != ' ' && data[i] != '\t' {
			return 0
		}
	}

	return i + 1
}
