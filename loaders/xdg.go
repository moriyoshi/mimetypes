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
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/moriyoshi/mimetypes"
)

func LoadXDGGlobsFile(r io.Reader) (mimetypes.MediaTypeRegistry, error) {
	mtr := mimetypes.NewBasicMediaTypeRegistry()

	b := bufio.NewScanner(r)
	b.Split(bufio.ScanLines)

	ln := 0
	for b.Scan() {
		ln += 1
		l := b.Text()
		l = strings.TrimFunc(l, unicode.IsSpace)
		if len(l) == 0 || l[0] == '#' {
			continue
		}

		fields := strings.Split(l, ":")
		if len(fields) != 2 || len(fields[1]) < 1 {
			return nil, fmt.Errorf("invalid entry at line %d", ln)
		}

		mtr.Add(
			mimetypes.MediaType{
				Name:  fields[0],
				Globs: []string{fields[1]},
			},
		)
	}

	if b.Err() != nil {
		return nil, b.Err()
	}

	return mtr, nil
}
