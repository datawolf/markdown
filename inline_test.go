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
	opts.Extensions |= EXTENSION_STRIKETHROUGH
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
	tests := []string{
		"some para_graph with _emphasised_ text.\n",
		"<p>some para_graph with <em>emphasised</em> text.</p>\n",

		"un*frigging*believable\n",
		"<p>un*frigging*believable</p>\n",
	}
	options := Options{Extensions: EXTENSION_NO_INTRA_EMPHASIS}
	doTestsInlineParam(t, tests, options, 0, HtmlRendererParameters{})
}

func TestStrong(t *testing.T) {
	tests := []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple **inline** test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"**at the** beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the **end**\n",
		"<p>at the <strong>end</strong></p>\n",

		"**try two** in **one line**\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over **two\nlines** test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd **number of** marker** here\n",
		"<p>odd <strong>number of</strong> marker** here</p>\n",

		"odd **number\nof** marker** here\n",
		"<p>odd <strong>number\nof</strong> marker** here</p>\n",

		"simple __inline__ test\n",
		"<p>simple <strong>inline</strong> test</p>\n",

		"__at the__ beginning\n",
		"<p><strong>at the</strong> beginning</p>\n",

		"at the __end__\n",
		"<p>at the <strong>end</strong></p>\n",

		"__try two__ in __one line__\n",
		"<p><strong>try two</strong> in <strong>one line</strong></p>\n",

		"over __two\nlines__ test\n",
		"<p>over <strong>two\nlines</strong> test</p>\n",

		"odd __number of__ marker__ here\n",
		"<p>odd <strong>number of</strong> marker__ here</p>\n",

		"odd __number\nof__ marker__ here\n",
		"<p>odd <strong>number\nof</strong> marker__ here</p>\n",

		"mix of **markers__\n",
		"<p>mix of **markers__</p>\n",

		"**`/usr`**: this folder is name `usr`\n",
		"<p><strong><code>/usr</code></strong>: this folder is name <code>usr</code></p>\n",

		"**`/usr`**:\n\n this folder is name `usr`\n",
		"<p><strong><code>/usr</code></strong>:</p>\n\n<p>this folder is name <code>usr</code></p>\n",
	}

	doTestsInline(t, tests)
}

func TestEmphasisMix(t *testing.T) {
	tests := []string{
		"***triple emphasis***\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple emphasis ***\n",
		"<p>***triple emphasis ***</p>\n",

		"***triple\nemphasis***\n",
		"<p><strong><em>triple\nemphasis</em></strong></p>\n",

		"___triple emphasis___\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"***triple emphasis___\n",
		"<p>***triple emphasis___</p>\n",

		"*__triple emphasis__*\n",
		"<p><em><strong>triple emphasis</strong></em></p>\n",

		"__*triple emphasis*__\n",
		"<p><strong><em>triple emphasis</em></strong></p>\n",

		"**improper  *nesting** is* bad\n",
		"<p><strong>improper  *nesting</strong> is* bad</p>\n",

		"*improper  **nesting* is** bad\n",
		"<p>*improper  <strong>nesting* is</strong> bad</p>\n",
	}
	doTestsInline(t, tests)
}

func TestStrikeTrough(t *testing.T) {
	tests := []string{
		"nothing inline\n",
		"<p>nothing inline</p>\n",

		"simple ~~inline~~ test\n",
		"<p>simple <del>inline</del> test</p>\n",

		"~~at the~~ beginning\n",
		"<p><del>at the</del> beginning</p>\n",

		"at ~~the end~~\n",
		"<p>at <del>the end</del></p>\n",

		"~~try two~~ in ~~one line~~\n",
		"<p><del>try two</del> in <del>one line</del></p>\n",

		"over ~~two\nlines~~ test\n",
		"<p>over <del>two\nlines</del> test</p>\n",

		"odd ~~number of~~ markers~~ here\n",
		"<p>odd <del>number of</del> markers~~ here</p>\n",

		"odd ~~number\nof~~ markers~~ here\n",
		"<p>odd <del>number\nof</del> markers~~ here</p>\n",
	}

	doTestsInline(t, tests)
}

func TestCodeSpan(t *testing.T) {
	tests := []string{
		"`source code`\n",
		"<p><code>source code</code></p>\n",

		"` source code with spaces `\n",
		"<p><code>source code with spaces</code></p>\n",

		"a `single marker\n",
		"<p>a `single marker</p>\n",

		"a single multi-tick marker with ``` no text\n",
		"<p>a single multi-tick marker with ``` no text</p>\n",

		"makers with ` ` a space\n",
		"<p>makers with  a space</p>\n",

		"`source code` and a `stray\n",
		"<p><code>source code</code> and a `stray</p>\n",

		"`source *with* _wakward characters_ in it`\n",
		"<p><code>source *with* _wakward characters_ in it</code></p>\n",

		"`spilt over\ntwo lines`\n",
		"<p><code>spilt over\ntwo lines</code></p>\n",

		"```multiple ticks``` for the marker\n",
		"<p><code>multiple ticks</code> for the marker</p>\n",

		"```multiple ticks `with` ticks inside```\n",
		"<p><code>multiple ticks `with` ticks inside</code></p>\n",
	}
	doTestsInline(t, tests)
}

func TestLineBreak(t *testing.T) {
	tests := []string{
		"this line  \nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br />\nextra space</p>\n",
	}
	doTestsInline(t, tests)

	tests = []string{
		"this line  \nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line \ndoes not\n",
		"<p>this line\ndoes not</p>\n",

		"this line\\\nhas a break\n",
		"<p>this line<br />\nhas a break</p>\n",

		"this line\\ \ndoes not\n",
		"<p>this line\\\ndoes not</p>\n",

		"this has an   \nextra space\n",
		"<p>this has an<br />\nextra space</p>\n",
	}

	opts := Options{
		Extensions: EXTENSION_BACKSLASH_LINE_BREAK,
	}
	doTestsInlineParam(t, tests, opts, 0, HtmlRendererParameters{})
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
