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

package mimetypes

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// MediaType represents a single media type designation
type MediaType struct {
	Name  string
	Globs []string
}

type InternalMediaType struct {
	MediaType
	Extensions []string
}

// BasicMediaTypeRegistry provides the basic implentation of MediaTypeRegistry.
type BasicMediaTypeRegistry struct {
	ExtToType      map[string]*InternalMediaType
	ExactToType    map[string]*InternalMediaType
	PatternsToType []*InternalMediaType
	NameToType     map[string]*InternalMediaType
}

const metas = "*?[\\"

// Add() adds a new MediaType to the registry.
func (mtr *BasicMediaTypeRegistry) Add(mt MediaType) {
	imt, ok := mtr.NameToType[mt.Name]
	if !ok {
		imt = &InternalMediaType{
			MediaType: MediaType{
				Name:  mt.Name,
				Globs: []string{},
			},
			Extensions: []string{},
		}
		mtr.NameToType[mt.Name] = imt
	}

	for _, glob := range mt.Globs {
		if strings.IndexAny(glob, metas) >= 0 {
			if len(glob) > 2 && glob[0] == '*' && glob[1] == '.' && strings.IndexAny(glob[2:], metas) == -1 {
				imt.Extensions = append(
					imt.Extensions,
					glob[1:],
				)
				mtr.ExtToType[glob[1:]] = imt
			}
		} else {
			mtr.ExactToType[glob] = imt
		}
		mtr.PatternsToType = append(mtr.PatternsToType, imt)
	}

	imt.Globs = append(imt.Globs, mt.Globs...)
}

func (mtr *BasicMediaTypeRegistry) ExtensionsByType(typ string) ([]string, error) {
	imt, ok := mtr.NameToType[typ]
	if !ok {
		return nil, nil
	}
	return imt.Extensions, nil
}

func (mtr *BasicMediaTypeRegistry) TypeByExtension(ext string) string {
	imt, ok := mtr.ExtToType[ext]
	if !ok {
		return ""
	}
	return imt.Name
}

func (mtr *BasicMediaTypeRegistry) TypeByFilename(name string) string {
	var imt *InternalMediaType
	var ok bool
	var ext string
	imt, ok = mtr.ExactToType[name]
	if ok {
		goto found
	}

	ext = path.Ext(name)
	if ext != "" {
		imt, ok = mtr.ExtToType[ext]
		if ok {
			goto found
		}
	}

	for _, imt = range mtr.PatternsToType {
		for _, glob := range imt.Globs {
			ok, err := path.Match(glob, name)
			if err != nil {
				continue
			}
			if ok {
				goto found
			}
		}
	}

	return ""
found:
	return imt.Name
}

func NewBasicMediaTypeRegistry() *BasicMediaTypeRegistry {
	return &BasicMediaTypeRegistry{
		ExtToType:      map[string]*InternalMediaType{},
		ExactToType:    map[string]*InternalMediaType{},
		PatternsToType: []*InternalMediaType{},
		NameToType:     map[string]*InternalMediaType{},
	}
}

type MediaTypeRegistry interface {
	ExtensionsByType(string) ([]string, error)
	TypeByExtension(string) string
	TypeByFilename(string) string
}

type Loader func(io.Reader) (MediaTypeRegistry, error)

var loaders = map[string]Loader{}

func AddLoader(format string, loader Loader) {
	loaders[format] = loader
}

func Load(path, format string) (MediaTypeRegistry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	loader, ok := loaders[format]
	if !ok {
		return nil, fmt.Errorf("no such loader: %s", format)
	}
	return loader(f)
}
