//
// util_test.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

package markdown

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSanitizedString(t *testing.T) {
	assert := require.New(t)

	tests := map[string]string{
		"This is a header":             "this-is-a-header",
		"This is also      a header":   "this-is-also-a-header",
		"main.go":                      "main-go",
		"Acticle 123":                  "acticle-123",
		"<- Let's try this, shall we?": "let-s-try-this-shall-we",
		"        ":                     "",
		"Hello, 世界":                    "hello-世界",
	}

	for key, value := range tests {
		assert.Equal(SanitizedString(key), value)
	}
}
