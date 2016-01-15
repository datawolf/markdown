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
