// Copyright (c) 2019 Moriyoshi Koizumi
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package loaders

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadXDGGlobsFileOk(t *testing.T) {
	mtr, err := LoadXDGGlobsFile(
		strings.NewReader(`
# comment
# aaa
	# comment with leading spaces
application/x-foo:*.ext1
application/x-foo:*.ext2
application/x-bar:*.ext3
application/x-exact-foo:exact1
application/x-exact-foo:exact2
application/x-exact-bar:exact3
application/x-complex-glob-foo:com*ple*x
application/x-complex-glob-foo:c*omp*lex
application/x-complex-glob-bar:co*mple*x
`,
		),
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	cases := []struct {
		expected string
		name     string
	}{
		{"application/x-foo", "a.ext1"},
		{"application/x-foo", "aa.ext1"},
		{"application/x-foo", "a.ext2"},
		{"application/x-foo", "aa.ext2"},
		{"application/x-bar", "a.ext3"},
		{"application/x-bar", "aa.ext3"},
		{"application/x-exact-foo", "exact1"},
		{"application/x-exact-foo", "exact2"},
		{"application/x-exact-bar", "exact3"},
		{"application/x-complex-glob-foo", "complex"},
		{"application/x-complex-glob-foo", "commpleex"},
		{"application/x-complex-glob-foo", "ccompppplex"},
		{"application/x-complex-glob-bar", "coooompleeeeex"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s / %s", c.expected, c.name), func(t *testing.T) {
			ext := path.Ext(c.name)
			if ext != "" {
				assert.Equal(t, c.expected, mtr.TypeByExtension(ext))
			}
			assert.Equal(t, c.expected, mtr.TypeByFilename(c.name))
		})
	}
}
