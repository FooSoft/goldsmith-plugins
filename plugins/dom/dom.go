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
	return &domPlugin{callback}
}

type domPlugin struct {
	callback Processor
}

func (*domPlugin) Name() string {
	return "dom"
}

func (*domPlugin) Initialize(context goldsmith.Context) ([]goldsmith.Filter, error) {
	return []goldsmith.Filter{extension.New(".html", ".htm")}, nil
}

func (plugin *domPlugin) Process(context goldsmith.Context, f goldsmith.File) error {
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return err
	}

	if err := plugin.callback(doc); err != nil {
		return err
	}

	html, err := doc.Html()
	if err != nil {
		return err
	}

	nf := goldsmith.NewFileFromData(f.Path(), []byte(html), f.ModTime())
	nf.InheritValues(f)
	context.DispatchFile(nf)
	return nil
}
