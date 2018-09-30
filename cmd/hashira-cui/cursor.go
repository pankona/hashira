package main

type cursor struct {
	index       int
	focusedPane *Pane
}

func (c *cursor) sanitize(maxLen int) *cursor {
	ret := c

	if c.index < 0 {
		ret.index = 0
	} else if c.index > maxLen {
		ret.index = maxLen
	}

	if c.index > len(c.focusedPane.priorities)-1 {
		ret.index = len(c.focusedPane.priorities) - 1
	}

	return ret
}

func (c *cursor) moveUp() {
	c.index--
	c = c.sanitize(len(c.focusedPane.priorities))
}

func (c *cursor) moveDown() {
	c.index++
	c = c.sanitize(len(c.focusedPane.priorities))
}

func (c *cursor) changeFocusedPane(p *Pane) {
	c.focusedPane = p
	c = c.sanitize(len(c.focusedPane.priorities))
}
