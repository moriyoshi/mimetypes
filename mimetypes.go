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
)

// MediaType represents a single media type designation
type MediaType struct {
	Name       string
	Extensions []string
}

// BasicMediaTypeRegistry provides the basic implentation of MediaTypeRegistry.
type BasicMediaTypeRegistry struct {
	ExtToType  map[string]*MediaType
	NameToType map[string]*MediaType
}

// Add() adds a new MediaType to the registry.
func (mtr *BasicMediaTypeRegistry) Add(mt MediaType) {
	existingType, ok := mtr.NameToType[mt.Name]
	if ok {
		existingType.Extensions = append(
			append(
				make([]string, 0, len(existingType.Extensions)+len(mt.Extensions)),
				existingType.Extensions...,
			),
			mt.Extensions...,
		)
		for _, ext := range mt.Extensions {
			mtr.ExtToType[ext] = existingType
		}
	} else {
		for _, ext := range mt.Extensions {
			mtr.ExtToType[ext] = &mt
		}
		mtr.NameToType[mt.Name] = &mt
	}
}

func (mtr *BasicMediaTypeRegistry) ExtensionsByType(typ string) ([]string, error) {
	mt, ok := mtr.NameToType[typ]
	if !ok {
		return nil, nil
	}
	return mt.Extensions, nil
}

func (mtr *BasicMediaTypeRegistry) TypeByExtension(ext string) string {
	mt, ok := mtr.ExtToType[ext]
	if !ok {
		return ""
	}
	return mt.Name
}

func NewBasicMediaTypeRegistry() *BasicMediaTypeRegistry {
	return &BasicMediaTypeRegistry{
		ExtToType:  map[string]*MediaType{},
		NameToType: map[string]*MediaType{},
	}
}

type MediaTypeRegistry interface {
	ExtensionsByType(string) ([]string, error)
	TypeByExtension(string) string
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
