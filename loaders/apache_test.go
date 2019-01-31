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
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadAPachestyleMimeTypeFileOk(t *testing.T) {
	mtr, err := LoadApacheStyleMimeTypeFile(
		strings.NewReader(`
# comment
# aaa
	# comment with leading spaces

application/x-foo	ext1 ext2
application/x-bar	.ext3
`,
		),
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, "application/x-foo", mtr.TypeByExtension(".ext1"))
	assert.Equal(t, "application/x-foo", mtr.TypeByExtension(".ext2"))
	assert.Equal(t, "application/x-bar", mtr.TypeByExtension(".ext3"))
}
