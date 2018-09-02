package main

import (
	"fmt"

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

		// render tasks for this pane
		for i, task := range p.tasks {
			if task.Place == p.place && !task.IsDeleted {
				var err error
				if i == selectedIndex {
					_, err = fmt.Fprintf(v, "\033[3%d;%dm%s\033[0m\n", 7, 4, task.Name)
				} else {
					_, err = fmt.Fprintf(v, "%s\n", task.Name)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
