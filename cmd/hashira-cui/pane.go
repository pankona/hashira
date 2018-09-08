package main

import (
	"fmt"
	"io"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

type Pane struct {
	name  string
	index int // place of this pane
	left  *Pane
	right *Pane
	place service.Place
	tasks []*service.Task
}

func (p *Pane) Layout(g *gocui.Gui, selectedIndex int) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(p.name, maxX/4*p.index, 1, maxX/4*p.index+maxX/4-1, maxY-1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
	}

	v.Clear()
	return renderTasks(v, p.tasks, selectedIndex)
}

func (p *Pane) len() int {
	return len(p.tasks)
}

func renderTasks(w io.Writer, tasks []*service.Task, selectedIndex int) error {
	var itemNum int
	// render tasks for this pane
	for _, task := range tasks {
		var err error
		if itemNum == selectedIndex {
			_, err = fmt.Fprintf(w, "\033[3%d;%dm%s\033[0m\n", 7, 4, task.Name)
		} else {
			_, err = fmt.Fprintf(w, "%s\n", task.Name)
		}
		if err != nil {
			return err
		}
		itemNum++
	}
	return nil
}
