package main

import (
	"github.com/jroimartin/gocui"
)

func configureKeyBindnig(g *gocui.Gui) error {
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
}

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}
