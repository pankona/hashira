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

func (p *Pane) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(p.name, maxX/4*p.index, 1, maxX/4*p.index+maxX/4-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
		for _, task := range p.tasks {
			if task.Place == p.place && !task.IsDeleted {
				_, err = fmt.Fprintln(v, task.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
