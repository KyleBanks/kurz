package main

import (
	"fmt"
	"os"

	"github.com/KyleBanks/kurz/pkg/doc"
	"github.com/KyleBanks/kurz/pkg/doc/parser"
	"github.com/KyleBanks/kurz/pkg/doc/resolver"
	"github.com/KyleBanks/kurz/pkg/ui"
	"github.com/KyleBanks/kurz/pkg/ui/console"
)

var path string

func init() {
	if len(os.Args) < 2 {
		printUsage(1)
	}

	switch arg := os.Args[1]; arg {

	case "-h":
		fallthrough
	case "--help":
		printUsage(0)

	default:
		path = arg

	}
}

func main() {
	r := resolver.Chain{
		Resolvers: []doc.Resolver{
			resolver.File{},
			resolver.URL{},
			resolver.Git{},
		},
	}
	p := parser.NewMarkdown(console.Styler{})

	runWithTUI(r, p)
}

func runWithTUI(r doc.Resolver, p doc.Parser) {
	w := console.NewWindow()
	w.ShowMessage(fmt.Sprintf("Loading %v...", path))

	go render(w, path, r, p, logError)

	if err := w.Run(); err != nil {
		logError(err)
	}
}

func render(c ui.Canvas, path string, r doc.Resolver, p doc.Parser, onErr func(error)) {
	d, err := doc.NewDocument(path, r, p)
	if err != nil {
		onErr(err)
		return
	}

	c.RenderDocument(d)
}
