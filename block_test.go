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
