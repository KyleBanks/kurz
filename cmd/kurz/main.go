package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KyleBanks/kurz/pkg/doc"
	"github.com/KyleBanks/kurz/pkg/doc/parser"
	"github.com/KyleBanks/kurz/pkg/doc/resolver"
	"github.com/KyleBanks/kurz/pkg/ui"
	"github.com/KyleBanks/kurz/pkg/ui/console"
	"github.com/KyleBanks/kurz/pkg/ui/text"
)

var (
	path string
	raw  bool
)

func init() {
	flag.Usage = printUsage

	flag.BoolVar(&raw, "raw", false, "If set, prints the markdown document as plain text.")
	flag.Parse()

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}
	path = os.Args[len(os.Args)-1]
}

func main() {
	r := resolver.Chain{
		Resolvers: []doc.Resolver{
			resolver.File{},
			resolver.URL{},
			resolver.Git{},
		},
	}
	var p parser.Markdown

	if raw {
		var t text.Renderer
		loadDoc(t, path, r, p, logError)
	} else {
		w := console.NewWindow()
		w.ShowMessage(fmt.Sprintf("Loading %v...", path))

		go loadDoc(w, path, r, p, logError)

		if err := w.Run(); err != nil {
			logError(err)
		}
	}
}

func loadDoc(c ui.Canvas, path string, r doc.Resolver, p doc.Parser, onErr func(error)) {
	d, err := doc.NewDocument(path, r, p)
	if err != nil {
		onErr(err)
		return
	}

	c.RenderDocument(d)
}
