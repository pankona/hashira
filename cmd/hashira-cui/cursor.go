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

	if c.index > c.focusedPane.tasks.Len()-1 {
		ret.index = c.focusedPane.tasks.Len() - 1
	}

	return ret
}
