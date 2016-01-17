//
// inline_test.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func doTestsInline(t *testing.T, tests []string) {
	doTestsInlineParam(t, tests, Options{}, 0, HtmlRendererParameters{})
}

func doTestsInlineParam(t *testing.T, tests []string, opts Options, htmlFlags int,
	params HtmlRendererParameters) {
	var candidate string

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdownInline(candidate, opts, htmlFlags, params)
		if actual != expected {
			t.Errorf("\n Input	[%#v]\nExpected[%#v]\nActual	[%#v]",
				candidate, expected, actual)
		}
	}
}

func runMarkdownInline(input string, opts Options, htmlFlags int, params HtmlRendererParameters) string {
	htmlFlags |= HTML_USE_XHTML

	renderer := HtmlRendererWithParameters(htmlFlags, "", "", params)

	return string(MarkdownOptions([]byte(input), renderer, opts))
}

//
//
// TestCases
//
//
func TestEmphasis(t *testing.T) {
	var tests = []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple *inline* test\n",
		"<p>simple <em>inline</em> test</p>\n",

		"*at the* beginning\n",
		"<p><em>at the</em> beginning</p>\n",

		"at the *end*\n",
		"<p>at the <em>end</em></p>\n",

		"*try two* in *one line*\n",
		"<p><em>try two</em> in <em>one line</em></p>\n",

		"over *tow\nlines* test\n",
		"<p>over <em>tow\nlines</em> test</p>\n",

		"odd *number of* markers* here\n",
		"<p>odd <em>number of</em> markers* here</p>\n",

		"odd *number\nof* markers* here\n",
		"<p>odd <em>number\nof</em> markers* here</p>\n",

		"simple _inline_ test\n",
		"<p>simple <em>inline</em> test</p>\n",

		"_at the_ beginning\n",
		"<p><em>at the</em> beginning</p>\n",

		"at the _end_\n",
		"<p>at the <em>end</em></p>\n",

		"_try two_ in _one line_\n",
		"<p><em>try two</em> in <em>one line</em></p>\n",

		"over _tow\nlines_ test\n",
		"<p>over <em>tow\nlines</em> test</p>\n",

		"odd _number of_ markers_ here\n",
		"<p>odd <em>number of</em> markers_ here</p>\n",

		"odd _number\nof_ markers_ here\n",
		"<p>odd <em>number\nof</em> markers_ here</p>\n",

		"mix of *markers_\n",
		"<p>mix of *markers_</p>\n",

		"*What is A\\* algorithm?*\n",
		"<p><em>What is A* algorithm?</em></p>\n",

		"some para_graph with _emphasised_ text.\n",
		"<p>some para_graph with <em>emphasised</em> text.</p>\n",

		"some paragraph with _emphasised_ te_xt.\n",
		"<p>some paragraph with <em>emphasised</em> te_xt.</p>\n",

		"some paragraph with t_wo bi_ts of _emphasised_ text.\n",
		"<p>some paragraph with t<em>wo bi</em>ts of <em>emphasised</em> text.</p>\n",

		"un*frigging*believable\n",
		"<p>un<em>frigging</em>believable</p>\n",
	}

	doTestsInline(t, tests)
}

func TestNoIntraEmphasis(t *testing.T) {
	tests := []string {
		"some para_graph with _emphasised_ text.\n",
		"<p>some para_graph with <em>emphasised</em> text.</p>\n",

		"un*frigging*believable\n",
		"<p>un*frigging*believable</p>\n",
	}
	options := Options{Extensions: EXTENSION_NO_INTRA_EMPHASIS}
	doTestsInlineParam(t, tests, options, 0, HtmlRendererParameters{})
}

//
//
// Unit TestCases
//
//
type testHelpFindEmphChar struct {
	text string
	ch   byte
	ret  int
}

func TestHelperFindEmphChar(t *testing.T) {
	assert := require.New(t)

	tests := []testHelpFindEmphChar{
		{"emph*", byte('*'), 4},
		{"emph~~[]*\n", byte('*'), 8},
		{"emph_", byte('_'), 4},
		{"emph[]~~_\n", byte('_'), 8},
		{"emp\\*h[]~~*\n", byte('*'), 10},
	}

	for _, test := range tests {
		assert.Equal(helperFindEmphChar([]byte(test.text), test.ch), test.ret)
	}
}
