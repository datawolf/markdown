//
// inline.go
// Copyright (C) 2016 wanglong <wanglong@wanglong-Lenovo-Product>
//
// Distributed under terms of the MIT license.
//

//
// Functions to parse inline elements

package markdown

import (
	"bytes"
)

// Functons to parse text with a block
// Each function returns the number of chars taken care of.
// input: is the complete block being rendererd
// offset: is the number of valid chars before the current cursor

func (p *parser) inline(out *bytes.Buffer, input []byte) {
	// this is called recurively: enforce a maximum depth
	if p.nesting >= p.maxNesting {
		return
	}

	p.nesting++

	i, end := 0, 0

	for i < len(input) {
		// copy inactive chars into output
		for end < len(input) && p.inlineCallback[input[end]] == nil {
			end++
		}

		p.r.NormalText(out, input[i:end])

		if end >= len(input) {
			break
		}

		i = end

		// call the grigger
		handler := p.inlineCallback[input[end]]
		if consumed := handler(p, out, input, i); consumed == 0 {
			// no action from the callback; buffer the byte for later
			end = i + 1
		} else {
			// skip past whatever the callback used
			i += consumed
			end = i
		}
	}
	p.nesting--
}

// `\\` backslash escape
var escapeChars = []byte("\\`*_{}[]()#+-.!:|&<>~")

func escape(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	data = data[offset:]

	if len(data) > 1 {
		if bytes.IndexByte(escapeChars, data[1]) < 0 {
			return 0
		}
		p.r.NormalText(out, data[1:2])
	}
	return 2
}

// single and double emphasis parsing
func emphasis(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	data = data[offset:]
	c := data[0]
	ret := 0

	// process: *test*  _test_
	if len(data) > 2 && data[1] != c {
		if isspace(data[1]) {
			return 0
		}
		if ret = helperEmphasis(p, out, data[1:], c); ret == 0 {
			return 0
		}
		return ret + 1
	}
	// process: **test**  __test__
	if len(data) > 3 && data[1] == c && data[2] != c {
		if isspace(data[2]) {
			return 0
		}
		if ret = helperDoubleEmphasis(p, out, data[2:], c); ret == 0 {
			return 0
		}
		return ret + 2
	}

	// process: ***test***  ___test___
	if len(data) > 4 && data[1] == c && data[2] == c && data[3] != c {
		if isspace(data[3]) {
			return 0
		}
		if ret = helperTripleEmphasis(p, out, data, 3, c); ret == 0 {
			return 0
		}

		return ret + 3
	}
	return 0
}

// helpFindEmphChar look for the next emph char, skipping other constructs
func helperFindEmphChar(data []byte, c byte) int {
	i := 0

	for i < len(data) {
		for i < len(data) && data[i] != c {
			i++
		}

		if i >= len(data) {
			return 0
		}
		// do not count escaped chars
		if i != 0 && data[i-1] == '\\' {
			i++
			continue
		}

		if data[i] == c {
			return i
		}
	}

	return 0
}
func helperEmphasis(p *parser, out *bytes.Buffer, data []byte, c byte) int {
	i := 0

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length
		if i >= len(data) {
			return 0
		}

		if i+1 < len(data) && data[i+1] == c {
			i++
			continue
		}
		if data[i] == c && !isspace(data[i-1]) {
			if p.flags&EXTENSION_NO_INTRA_EMPHASIS != 0 {
				if !(i+1 == len(data) || isspace(data[i+1]) || ispunct(data[i+1])) {
					continue
				}
			}
			var work bytes.Buffer
			p.inline(&work, data[:i])
			p.r.Emphasis(out, work.Bytes())
			return i + 1
		}
	}
	return 0
}

func helperDoubleEmphasis(p *parser, out *bytes.Buffer, data []byte, c byte) int {
	i := 0

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		if i+1 < len(data) && data[i] == c && data[i+1] == c && i > 0 && !isspace(data[i-1]) {
			var work bytes.Buffer
			p.inline(&work, data[:i])

			if work.Len() > 0 {
				p.r.DoubleEmphasis(out, work.Bytes())
			}

			return i + 2
		}
		i++
	}
	return 0
}

func helperTripleEmphasis(p *parser, out *bytes.Buffer, data []byte, offset int, c byte) int {
	i := 0
	origData := data
	data = data[offset:]

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		// skip whitespace  proceded symbols
		if data[i] != c || isspace(data[i-1]) {
			continue
		}

		switch {
		case i+2 < len(data) && data[i+1] == c && data[i+2] == c:
			// triple symbol found
			var work bytes.Buffer

			p.inline(&work, data[:i])
			if work.Len() > 0 {
				p.r.TripleEmphasis(out, work.Bytes())
			}
			return i + 3
		case i+1 < len(data) && data[i+1] == c:
			// double symbol found, hand over to emph1
			length = helperEmphasis(p, out, origData[offset-2:], c)
			if length == 0 {
				return 0
			} else {
				return length - 2
			}
		default:
			// single symbol found, hand over to emph2
			length = helperDoubleEmphasis(p, out, origData[offset-1:], c)
			if length == 0 {
				return 0
			} else {
				return length - 1
			}
		}
	}
	return 0
}
