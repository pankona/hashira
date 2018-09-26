package main

import (
	"fmt"
	"io"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

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

func (p *Pane) Layout(g *gocui.Gui, c *cursor, focusedIndex int, selectedTask *service.Task) error {
	maxX, maxY := g.Size()
	r := rectangle{maxX / 4 * p.index, 1, maxX/4*p.index + maxX/4 - 1, maxY - 1}
	v, err := g.SetView(p.name, r.x0, r.y0, r.x1, r.y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
	}

	v.Clear()

	return p.renderTasks(v, r, c, p.tasks, p.priorities, focusedIndex, selectedTask)
}

func (p *Pane) len() int {
	return len(p.tasks)
}

func (p *Pane) renderTasks(w io.Writer,
	rect rectangle, c *cursor, tasks map[string]*service.Task, priorities []string,
	focusedIndex int, selectedTask *service.Task) error {
	var err error

	height := rect.y1 - rect.y0
	if height < 0 {
		return fmt.Errorf("invalid pane height. height must be positive")
	}

	// cursor must be in pane height
	if c.index < 0 {
		c.index = 0
	} else if c.index >= height-2 {
		c.index = height - 2
	}

	// calculate scroll
	to := p.renderFrom + height - 2 // -2 for considering frame width
	if focusedIndex == -1 {
		// this pane is not focused. nop
	} else if focusedIndex > to {
		p.renderFrom += focusedIndex - to
	} else if focusedIndex < p.renderFrom {
		p.renderFrom -= p.renderFrom - focusedIndex
	}
	to = p.renderFrom + height - 2

	if focusedIndex != -1 {
		log.Printf("[%s] focused index = %d, from:to = %d:%d, cursor index = %d",
			p.name, focusedIndex, p.renderFrom, to, c.index)
	}

	var taskNum int

	if p == c.pane && c.index > len(priorities)-1 {
		c.index = len(priorities) - 1
	}

	// render tasks for this pane
	for i, id := range priorities {
		if i < p.renderFrom || i > to {
			continue
		}

		task, ok := tasks[id]
		if !ok {
			// should not reach here
			// TODO: error logging and continue
			continue
		}

		prefix := ""
		if selectedTask != nil && task.Id == selectedTask.Id {
			prefix = "*"
		}

		//if task == focusedTask {
		if p == c.pane && taskNum == c.index {
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
