// Package dom makes it easy to modify your document structure.
package dom

import (
	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-components/filters/extension"
	"github.com/PuerkitoBio/goquery"
)

// Dom plugin context.
type Dom interface {
	goldsmith.Plugin
	goldsmith.Initializer
	goldsmith.Processor
}

// A Processor callback function is used to modify documents.
type Processor func(*goquery.Document) error

// New creates a new instance of the Dom plugin.
func New(callback Processor) Dom {
	return &dom{callback}
}

type dom struct {
	callback Processor
}

func (*dom) Name() string {
	return "dom"
}

func (*dom) Initialize(context *goldsmith.Context) ([]goldsmith.Filter, error) {
	return []goldsmith.Filter{extension.New(".html", ".htm")}, nil
}

func (d *dom) Process(context *goldsmith.Context, inputFile *goldsmith.File) error {
	if outputFile := context.RetrieveCachedFile(inputFile.Path(), inputFile); outputFile != nil {
		context.DispatchFile(outputFile)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(inputFile)
	if err != nil {
		return err
	}

	if err := d.callback(doc); err != nil {
		return err
	}

	html, err := doc.Html()
	if err != nil {
		return err
	}

	outputFile := goldsmith.NewFileFromData(inputFile.Path(), []byte(html))
	outputFile.InheritValues(inputFile)
	context.DispatchAndCacheFile(outputFile)

	return nil
}
