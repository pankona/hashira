package main

import (
	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

type Ctrl struct {
	hashirac *hashirac.Client
}

func NewCtrl(cli *hashirac.Client) *Ctrl {
	return &Ctrl{hashirac: cli}
}

func (c *Ctrl) ConfigureKeyBindings(g *gocui.Gui) error {
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.Quit)
}

func (c *Ctrl) Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}
