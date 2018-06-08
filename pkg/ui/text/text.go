package text

import (
	"fmt"

	"github.com/KyleBanks/kurz/pkg/doc"
)

type Renderer struct{}

func (Renderer) RenderDocument(d doc.Document) {
	for _, h := range d.Headings {
		fmt.Println(h.Title)
		for _, s := range h.Content {
			fmt.Println(s.Text)
		}
	}
}
