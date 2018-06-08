package ui

import (
	"github.com/KyleBanks/kurz/pkg/doc"
)

type Canvas interface {
	RenderDocument(doc.Document)
}
