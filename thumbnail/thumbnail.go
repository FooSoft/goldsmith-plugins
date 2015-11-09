/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package thumbnail

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"sync"

	"github.com/FooSoft/goldsmith"
)

type Namer func(path string, dims uint) (string, bool)

type thumbnail struct {
	dims  uint
	namer Namer
}

func (t *thumbnail) New(dims uint, namer Namer) (goldsmith.Chainer, error) {
	return &thumbnail{dims, namer}, nil
}

func (t *thumbnail) Filter(path string) bool {
	switch filepath.Ext(path) {
	case ".jpeg":
		fallthrough
	case ".jpg":
		fallthrough
	case ".gif":
		fallthrough
	case ".png":
		return true
	default:
		return false
	}
}

func (t *thumbnail) thumbnail(ctx goldsmith.Context, origFile *goldsmith.File, thumbPath string) (*goldsmith.File, error) {
	origImg, _, err := image.Decode(&origFile.Buff)
	if err != nil {
		return nil, err
	}

	thumbFile := ctx.NewFile(thumbPath)

	switch filepath.Ext(thumbPath) {
	case ".jpeg":
		fallthrough
	case ".jpg":
		thumbFile.Err = jpeg.Encode(&thumbFile.Buff, origImg, nil)
	case ".gif":
		thumbFile.Err = gif.Encode(&thumbFile.Buff, origImg, nil)
	case ".png":
		thumbFile.Err = png.Encode(&thumbFile.Buff, origImg)
	}

	return thumbFile, nil
}

func (t *thumbnail) Chain(ctx goldsmith.Context, input, output chan *goldsmith.File) {
	defer close(output)

	var wg sync.WaitGroup
	for file := range input {
		wg.Add(1)
		go func(f *goldsmith.File) {
			defer wg.Done()

			if path, create := t.namer(file.Path, t.dims); create {
				if thumb, err := t.thumbnail(ctx, f, path); err == nil {
					output <- thumb
				}
			}

			output <- f
			wg.Done()
		}(file)
	}

	wg.Wait()
}