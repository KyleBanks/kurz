package console

import "fmt"

const (
	collapsedContent = "..."
)

// contentState maps a heading + section index to a
// special state string.
type contentState struct {
	store map[string]string
}

// newContentState initializes and returns a new contentState type.
func newContentState() *contentState {
	return &contentState{
		store: make(map[string]string),
	}
}

// get returns the stored content state for a given heading+section index.
func (c *contentState) get(heading, section int) string {
	s, ok := c.store[contentKey(heading, section)]
	if !ok {
		return ""
	}

	return s
}

// set updates the stored content state for a given heading+section index.
func (c *contentState) set(heading, section int, s string) {
	c.store[contentKey(heading, section)] = s
}

// clear removes the stored content state for a given heading+section index.
func (c *contentState) clear(heading, section int) {
	delete(c.store, contentKey(heading, section))
}

// contentKey returns a unique key for any heading+section combination.
func contentKey(heading int, section int) string {
	return fmt.Sprintf("%d:%d", heading, section)
}
