package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type view struct {
	ps *PubSub
}

func (v *view) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err = fmt.Fprintln(v, "Hello world!")
		if err != nil {
			return err
		}
	}
	return nil
}
