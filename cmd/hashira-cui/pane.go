package main

import (
	"fmt"
	"io"

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
	from       int
}

type rectangle struct {
	x0, y0, x1, y1 int
}

func (p *Pane) Layout(g *gocui.Gui, focusedTask, selectedTask *service.Task) error {
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

	return p.renderTasks(v, r, p.tasks, p.priorities, focusedTask, selectedTask)
}

func (p *Pane) len() int {
	return len(p.tasks)
}

func (p *Pane) renderTasks(w io.Writer, rect rectangle, tasks map[string]*service.Task, priorities []string, focusedTask, selectedTask *service.Task) error {
	var taskNum int
	var err error

	height := rect.y1 - rect.y0
	if height < 0 {
		return fmt.Errorf("invalid pane height. height must be positive")
	}

	var focusedIndex int
	for i, id := range priorities {
		task, ok := tasks[id]
		if !ok {
			// should not reach here
			// TODO: error logging and continue
			continue
		}
		if task == focusedTask {
			focusedIndex = i
		}
	}

	to := p.from + height - 2 // -2 for frame width
	if focusedIndex > to {
		p.from += focusedIndex - to
	} else if focusedIndex < p.from {
		p.from -= p.from - focusedIndex
	}
	to = p.from + height - 2

	// render tasks for this pane
	for i, id := range priorities {
		if i < p.from || i > to {
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

		if task == focusedTask {
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
