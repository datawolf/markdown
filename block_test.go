//
// block_test.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func runMarkdownBlockWithRenderer(input string, extensions int, renderer Renderer) string {
	return string(Markdown([]byte(input), renderer, extensions))
}

func runMarkdownBlock(input string, extensions int) string {
	htmlFlags := 0
	htmlFlags |= HTML_USE_XHTML

	renderer := HtmlRenderer(htmlFlags, "", "")

	return runMarkdownBlockWithRenderer(input, extensions, renderer)
}

func runnerWithRendererParameters(parameters HtmlRendererParameters) func(string, int) string {
	return func(input string, extensions int) string {
		htmlFlags := 0
		htmlFlags |= HTML_USE_XHTML

		renderer := HtmlRendererWithParameters(htmlFlags, "", "", parameters)

		return runMarkdownBlockWithRenderer(input, extensions, renderer)
	}
}

func doTestsBlock(t *testing.T, tests []string, extensions int) {
	doTestsBlockWithRunner(t, tests, extensions, runMarkdownBlock)
}

func doTestsBlockWithRunner(t *testing.T, tests []string, extensions int, runner func(string, int) string) {
	// catch and report panics
	var candidate string
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("\npanic while processing [%#v]: %s\n", candidate, err)
		}
	}()

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runner(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput	[%#v]\nExpected[%#v]\nActual	[%#v]",
				candidate, expected, actual)
		}
	}
}

//
//
//
// The TestCases
//
//
//
func TestPrefixHeaderNoExtensions(t *testing.T) {
	tests := []string{
		"# Header 1\n",
		"<h1>Header 1</h1>\n",

		"## Header 2\n",
		"<h2>Header 2</h2>\n",

		"### Header 3\n",
		"<h3>Header 3</h3>\n",

		"#### Header 4\n",
		"<h4>Header 4</h4>\n",

		"##### Header 5\n",
		"<h5>Header 5</h5>\n",

		"###### Header 6\n",
		"<h6>Header 6</h6>\n",

		"####### Header 7\n",
		"<h6># Header 7</h6>\n",

		"#Header 1\n",
		"<h1>Header 1</h1>\n",

		"##Header 2\n",
		"<h2>Header 2</h2>\n",

		"###Header 3\n",
		"<h3>Header 3</h3>\n",

		"####Header 4\n",
		"<h4>Header 4</h4>\n",

		"#####Header 5\n",
		"<h5>Header 5</h5>\n",

		"######Header 6\n",
		"<h6>Header 6</h6>\n",

		"#######Header 7\n",
		"<h6>#Header 7</h6>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",

		"#Header 1 \\#\n",
		"<h1>Header 1 #</h1>\n",

		"#Header 1 \\# foo\n",
		"<h1>Header 1 # foo</h1>\n",

		"#Header 1 #\\##\n",
		"<h1>Header 1 ##</h1>\n",
	}

	doTestsBlock(t, tests, 0)
}

func TestPrefixHeaderSpaceExtension(t *testing.T) {
	tests := []string{
		"# Header 1\n",
		"<h1>Header 1</h1>\n",

		"## Header 2\n",
		"<h2>Header 2</h2>\n",

		"### Header 3\n",
		"<h3>Header 3</h3>\n",

		"#### Header 4\n",
		"<h4>Header 4</h4>\n",

		"##### Header 5\n",
		"<h5>Header 5</h5>\n",

		"###### Header 6\n",
		"<h6>Header 6</h6>\n",

		"####### Header 7\n",
		"<p>####### Header 7</p>\n",

		"#Header 1\n",
		"<p>#Header 1</p>\n",

		"##Header 2\n",
		"<p>##Header 2</p>\n",

		"###Header 3\n",
		"<p>###Header 3</p>\n",

		"####Header 4\n",
		"<p>####Header 4</p>\n",

		"#####Header 5\n",
		"<p>#####Header 5</p>\n",

		"######Header 6\n",
		"<p>######Header 6</p>\n",

		"#######Header 7\n",
		"<p>#######Header 7</p>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_SPACE_HEADERS)
}

func TestPrefixHeaderIdExtension(t *testing.T) {
	var tests = []string{
		"# Header 1 {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1 {#someid}   \n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1         {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1 {#someid\n",
		"<h1>Header 1 {#someid</h1>\n",

		"# Header 1 {#someid\n",
		"<h1>Header 1 {#someid</h1>\n",

		"# Header 1 {#someid}}\n",
		"<h1 id=\"someid\">Header 1</h1>\n\n<p>}</p>\n",

		"## Header 2 {#someid}\n",
		"<h2 id=\"someid\">Header 2</h2>\n",

		"### Header 3 {#someid}\n",
		"<h3 id=\"someid\">Header 3</h3>\n",

		"#### Header 4 {#someid}\n",
		"<h4 id=\"someid\">Header 4</h4>\n",

		"##### Header 5 {#someid}\n",
		"<h5 id=\"someid\">Header 5</h5>\n",

		"###### Header 6 {#someid}\n",
		"<h6 id=\"someid\">Header 6</h6>\n",

		"####### Header 7 {#someid}\n",
		"<h6 id=\"someid\"># Header 7</h6>\n",
		"# Header 1 # {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"## Header 2 ## {#someid}\n",
		"<h2 id=\"someid\">Header 2</h2>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_HEADER_IDS)
}

func TestPrefixHeaderIdExtensionWithPrefixAndSuffix(t *testing.T) {
	var tests = []string{
		"# header 1 {#someid}\n",
		"<h1 id=\"PRE:someid:POST\">header 1</h1>\n",

		"## header 2 {#someid}\n",
		"<h2 id=\"PRE:someid:POST\">header 2</h2>\n",

		"### header 3 {#someid}\n",
		"<h3 id=\"PRE:someid:POST\">header 3</h3>\n",

		"#### header 4 {#someid}\n",
		"<h4 id=\"PRE:someid:POST\">header 4</h4>\n",

		"##### header 5 {#someid}\n",
		"<h5 id=\"PRE:someid:POST\">header 5</h5>\n",

		"###### header 6 {#someid}\n",
		"<h6 id=\"PRE:someid:POST\">header 6</h6>\n",

		"####### header 7 {#someid}\n",
		"<h6 id=\"PRE:someid:POST\"># header 7</h6>\n",

		"# header 1 # {#someid}\n",
		"<h1 id=\"PRE:someid:POST\">header 1</h1>\n",

		"## header 2 ## {#someid}\n",
		"<h2 id=\"PRE:someid:POST\">header 2</h2>\n",
	}
	parameters := HtmlRendererParameters{
		HeaderIDPrefix: "PRE:",
		HeaderIDSuffix: ":POST",
	}

	doTestsBlockWithRunner(t, tests, EXTENSION_HEADER_IDS, runnerWithRendererParameters(parameters))
}

func TestPrefixAutoHeaderIdExtension(t *testing.T) {
	var tests = []string{
		"# Header 1\n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"# Header 1   \n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"## Header 2\n",
		"<h2 id=\"header-2\">Header 2</h2>\n",

		"### Header 3\n",
		"<h3 id=\"header-3\">Header 3</h3>\n",

		"#### Header 4\n",
		"<h4 id=\"header-4\">Header 4</h4>\n",

		"##### Header 5\n",
		"<h5 id=\"header-5\">Header 5</h5>\n",

		"###### Header 6\n",
		"<h6 id=\"header-6\">Header 6</h6>\n",

		"####### Header 7\n",
		"<h6 id=\"header-7\"># Header 7</h6>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1 id=\"header-1\">Header 1</h1>\n\n<p>Goodbye</p>\n",
		"# Header\n\n# Header\n",
		"<h1 id=\"header\">Header</h1>\n\n<h1 id=\"header-1\">Header</h1>\n",

		"# Header 1\n\n# Header 1",
		"<h1 id=\"header-1\">Header 1</h1>\n\n<h1 id=\"header-1-1\">Header 1</h1>\n",

		"# Header\n\n# Header 1\n\n# Header\n\n# Header",
		"<h1 id=\"header\">Header</h1>\n\n<h1 id=\"header-1\">Header 1</h1>\n\n<h1 id=\"header-1-1\">Header</h1>\n\n<h1 id=\"header-1-2\">Header</h1>\n",
	}
	doTestsBlock(t, tests, EXTENSION_AUTO_HEADER_IDS)
}

func TestPrefixAutoHeaderIdExtensionWithPrefixAndSuffix(t *testing.T) {
	var tests = []string{
		"# Header 1\n",
		"<h1 id=\"PRE:header-1:POST\">Header 1</h1>\n",

		"# Header 1   \n",
		"<h1 id=\"PRE:header-1:POST\">Header 1</h1>\n",

		"## Header 2\n",
		"<h2 id=\"PRE:header-2:POST\">Header 2</h2>\n",

		"### Header 3\n",
		"<h3 id=\"PRE:header-3:POST\">Header 3</h3>\n",

		"#### Header 4\n",
		"<h4 id=\"PRE:header-4:POST\">Header 4</h4>\n",

		"##### Header 5\n",
		"<h5 id=\"PRE:header-5:POST\">Header 5</h5>\n",

		"###### Header 6\n",
		"<h6 id=\"PRE:header-6:POST\">Header 6</h6>\n",

		"####### Header 7\n",
		"<h6 id=\"PRE:header-7:POST\"># Header 7</h6>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1 id=\"PRE:header-1:POST\">Header 1</h1>\n\n<p>Goodbye</p>\n",

		"# Header\n\n# Header\n",
		"<h1 id=\"PRE:header:POST\">Header</h1>\n\n<h1 id=\"PRE:header-1:POST\">Header</h1>\n",

		"# Header 1\n\n# Header 1",
		"<h1 id=\"PRE:header-1:POST\">Header 1</h1>\n\n<h1 id=\"PRE:header-1-1:POST\">Header 1</h1>\n",

		"# Header\n\n# Header 1\n\n# Header\n\n# Header",
		"<h1 id=\"PRE:header:POST\">Header</h1>\n\n<h1 id=\"PRE:header-1:POST\">Header 1</h1>\n\n<h1 id=\"PRE:header-1-1:POST\">Header</h1>\n\n<h1 id=\"PRE:header-1-2:POST\">Header</h1>\n",
	}

	parameters := HtmlRendererParameters{
		HeaderIDPrefix: "PRE:",
		HeaderIDSuffix: ":POST",
	}

	doTestsBlockWithRunner(t, tests, EXTENSION_AUTO_HEADER_IDS, runnerWithRendererParameters(parameters))
}

func TestPrefixMultipleHeaderExtensions(t *testing.T) {
	tests := []string{
		"# Header\n\n# Header {#header}\n\n# Header 1",
		"<h1 id=\"header\">Header</h1>\n\n<h1 id=\"header-1\">Header</h1>\n\n<h1 id=\"header-1-1\">Header 1</h1>\n",
	}
	doTestsBlock(t, tests, EXTENSION_AUTO_HEADER_IDS|EXTENSION_HEADER_IDS)
}

func TestUnderlineHeaders(t *testing.T) {
	tests := []string{
		"Header 1\n========\n",
		"<h1>Header 1</h1>\n",

		"Header 2\n--------\n",
		"<h2>Header 2</h2>\n",

		"A\n=\n",
		"<h1>A</h1>\n",

		"B\n-\n",
		"<h2>B</h2>\n",

		"Paragraph\nHeader\n=\n",
		"<p>Paragraph</p>\n\n<h1>Header</h1>\n",

		"Header\n===\nParagraph\n",
		"<h1>Header</h1>\n\n<p>Paragraph</p>\n",

		"Header\n===\nAnother header\n---\n",
		"<h1>Header</h1>\n\n<h2>Another header</h2>\n",

		"   Header\n======\n",
		"<h1>Header</h1>\n",

		"Header with *inline*\n=====\n",
		"<h1>Header with <em>inline</em></h1>\n",

		"Paragraph\n\n\n\n\nHeader\n===\n",
		"<p>Paragraph</p>\n\n<h1>Header</h1>\n",

		"Trailing space \n====        \n\n",
		"<h1>Trailing space</h1>\n",
		"Trailing spaces\n====        \n\n",
		"<h1>Trailing spaces</h1>\n",

		"Double underline\n=====\n=====\n",
		"<h1>Double underline</h1>\n\n<p>=====</p>\n",
	}

	doTestsBlock(t, tests, 0)
}

func TestUnderlineHeadersAutoIDs(t *testing.T) {
	var tests = []string{
		"Header 1\n========\n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"Header 2\n--------\n",
		"<h2 id=\"header-2\">Header 2</h2>\n",

		"A\n=\n",
		"<h1 id=\"a\">A</h1>\n",

		"B\n-\n",
		"<h2 id=\"b\">B</h2>\n",

		"Paragraph\nHeader\n=\n",
		"<p>Paragraph</p>\n\n<h1 id=\"header\">Header</h1>\n",

		"Header\n===\nParagraph\n",
		"<h1 id=\"header\">Header</h1>\n\n<p>Paragraph</p>\n",

		"Header\n===\nAnother header\n---\n",
		"<h1 id=\"header\">Header</h1>\n\n<h2 id=\"another-header\">Another header</h2>\n",

		"   Header\n======\n",
		"<h1 id=\"header\">Header</h1>\n",

		"Header with *inline*\n=====\n",
		"<h1 id=\"header-with-inline\">Header with <em>inline</em></h1>\n",

		"Paragraph\n\n\n\n\nHeader\n===\n",
		"<p>Paragraph</p>\n\n<h1 id=\"header\">Header</h1>\n",

		"Trailing space \n====        \n\n",
		"<h1 id=\"trailing-space\">Trailing space</h1>\n",

		"Trailing spaces\n====        \n\n",
		"<h1 id=\"trailing-spaces\">Trailing spaces</h1>\n",

		"Double underline\n=====\n=====\n",
		"<h1 id=\"double-underline\">Double underline</h1>\n\n<p>=====</p>\n",

		"Header\n======\n\nHeader\n======\n",
		"<h1 id=\"header\">Header</h1>\n\n<h1 id=\"header-1\">Header</h1>\n",

		"Header 1\n========\n\nHeader 1\n========\n",
		"<h1 id=\"header-1\">Header 1</h1>\n\n<h1 id=\"header-1-1\">Header 1</h1>\n",
	}
	doTestsBlock(t, tests, EXTENSION_AUTO_HEADER_IDS)
}

//
//
// Unit TestCases
//
//
func TestIsPrefixHeader(t *testing.T) {
	assert := require.New(t)

	tests := map[string]bool{
		"$":               false,
		"#header 1":       true,
		"# header 1":      true,
		"#######header 7": true,
	}

	p := new(parser)
	for key, value := range tests {
		assert.Equal(p.isPrefixHeader([]byte(key)), value)
	}

	testsExtSpaceHeader := map[string]bool{
		"$":                false,
		"# ":               true,
		"#header 1":        false,
		"####### header 7": false,
		"# header 1":       true,
		"### header 3":     true,
		"# he":             true,
	}
	p1 := new(parser)
	p1.flags |= EXTENSION_SPACE_HEADERS
	for key, value := range testsExtSpaceHeader {
		assert.Equal(p1.isPrefixHeader([]byte(key)), value)
	}
}

func TestIsEmpty(t *testing.T) {
	assert := require.New(t)

	tests := map[string]int{
		"":                 0,
		"\n":               1,
		"empty":            0,
		"empty\n":          0,
		"\t a   \nempty\n": 0,
		"\t    \nempty\n":  6,
	}
	p := new(parser)
	for key, value := range tests {
		assert.Equal(p.isEmpty([]byte(key)), value)
	}
}
