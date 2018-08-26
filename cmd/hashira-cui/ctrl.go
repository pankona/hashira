package main

import (
	"github.com/jroimartin/gocui"
)

func ConfigureKeyBindnig(g *gocui.Gui) error {
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
}

func quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}
