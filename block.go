//
// block.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"bytes"
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

		// anything else must look like a normal paragraph
		input = input[p.paragraph(out, input):]
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

	// get the id
	if p.flags&EXTENSION_HEADER_IDS != 0 {
		j, k := 0, 0

		// find the start/end of header id
		for j = start; j < end-1 && (input[j] != '{' || input[j+1] != '#'); j++ {
		}
		for k = j + 1; k < end && input[k] != '}'; k++ {
		}

		// extract the header id if found
		if j < end && k < end {
			id = string(input[j+2 : k])
			end = j
			skip = k + 1
			for end > 0 && input[end-1] == ' ' {
				end--
			}
		}
	}

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

func (p *parser) paragraph(out *bytes.Buffer, data []byte) int {
	var i int

	// keep going until we find something to mark the end of the paragraph
	for i < len(data) {
		// mark the beginning of the current line
		//	prev = line
		current := data[i:]

		// did we find a blank line marking the end of the paragraph
		if n := p.isEmpty(current); n > 0 {
			p.renderParagraph(out, data[:i])
			return i + n
		}

		// an underline under some text marks a header, so our paragraph ended on prev line

		// if there's a prefixed header paragraph is over
		if p.isPrefixHeader(current) {
			p.renderParagraph(out, data[:i])
			return i
		}

		// otherwise, scan to the beginning of the next line
		for data[i] != '\n' {
			i++
		}
		i++
	}

	p.renderParagraph(out, data[:i])
	return i
}

// renderParagraph render a single a paragraph that has already been parsed out
func (p *parser) renderParagraph(out *bytes.Buffer, data []byte) {
	if len(data) == 0 {
		return
	}

	// trim leading spaces
	begin := 0
	for data[begin] == ' ' {
		begin++
	}

	// trim trailing newline
	end := len(data) - 1

	// trim trailing spaces
	for end > begin && data[end-1] == ' ' {
		end--
	}

	work := func() bool {
		p.inline(out, data[begin:end])
		return true
	}

	p.r.Paragraph(out, work)
}
