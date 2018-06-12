package console

import "fmt"

var (
	commandsTableOfContents = []command{
		{" ESC ", "Exit"},
		{"⬆ ", "Up"},
		{"⬇ ", "Down"},
		{" ➡ / ENTER ", "Select"},
	}
	commandsContent = []command{
		{" ⬅ / ESC ", "Go Back"},
		{"⬆ ", "Up"},
		{"⬇ ", "Down"},
		{" SPACE ", "Collapse"},
	}
)

type command struct {
	symbol string
	label  string
}

func (c command) String() string {
	return fmt.Sprintf("[black:white:b]%v[-:-:-] %v", c.symbol, c.label)
}
