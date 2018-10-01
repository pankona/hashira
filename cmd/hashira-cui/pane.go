package main

import (
	"fmt"
	"io"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

// Pane represents a pane,
// like one of Backlog, ToDo, Doing, Done
type Pane struct {
	name       string
	index      int // place of this pane
	left       *Pane
	right      *Pane
	place      service.Place
	tasks      map[string]*service.Task
	priorities []string // array of task IDs
	renderFrom int
}

type rectangle struct {
	x0, y0, x1, y1 int
}

// Layout writes tasks in pane
func (p *Pane) Layout(g *gocui.Gui, c *cursor, focusedIndex int, selectedTask *service.Task) error {
	maxX, maxY := g.Size()
	rect := rectangle{maxX / 4 * p.index, 1, maxX/4*p.index + maxX/4 - 1, maxY - 1}

	v, err := g.SetView(p.name, rect.x0, rect.y0, rect.x1, rect.y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
	}

	v.Clear()

	return p.render(v, rect, c, focusedIndex, selectedTask)
}

func (p *Pane) render(w io.Writer, rect rectangle, cursor *cursor,
	focusedIndex int, selectedTask *service.Task) error {

	// -1 is adjustment for considering width of frame
	maxLen := rect.y1 - rect.y0 - 2
	if maxLen < 0 {
		return fmt.Errorf("invalid pane height. height must be positive")
	}

	// cursor must point within pane max length
	c := cursor.sanitize(maxLen)

	// calculate index from where to render for scrolling
	p.renderFrom = p.calcRenderFrom(focusedIndex, maxLen)

	return p.renderTasks(w, c, selectedTask)
}

func (p *Pane) calcRenderFrom(focusedIndex, maxLen int) int {
	renderFrom := p.renderFrom

	if focusedIndex == -1 {
		// this pane is not focused. nop
		return renderFrom
	}

	// calculate renderFrom for scrolling
	to := renderFrom + maxLen

	if focusedIndex > to {
		renderFrom += focusedIndex - to
	} else if focusedIndex < p.renderFrom {
		renderFrom -= renderFrom - focusedIndex
	}

	return renderFrom
}

func (p *Pane) renderTasks(w io.Writer, cursor *cursor, selected *service.Task) error {
	var taskNum int

	// render tasks for this pane
	for i, id := range p.priorities {
		if i < p.renderFrom {
			// skip rendering to scroll
			continue
		}

		task, ok := p.tasks[id]
		if !ok {
			// should not reach here
			// TODO: error logging and continue
			log.Printf("[WARNING] a task with ID [%s] is missing", id)
			continue
		}

		prefix := ""
		if selected != nil && task.Id == selected.Id {
			prefix = "*"
		}

		var err error

		if p == cursor.focusedPane && taskNum == cursor.index {
			_, err = fmt.Fprintf(w, "%s \033[3%d;%dm%s\033[0m\n", prefix, 7, 4, task.Name)
		} else {
			_, err = fmt.Fprintf(w, "%s %s\n", prefix, task.Name)
		}

		if err != nil {
			return err
		}

		taskNum++
	}

	return nil
}
