package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type View struct {
}

func (v *View) Layout(g *gocui.Gui) error {
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

func (v *View) OnEvent(event string, data interface{}) {

}
