package main

import (
	"context"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
	"github.com/pkg/errors"
)

type Ctrl struct {
	hashirac *hashirac.Client
	pub      Publisher
}

func NewCtrl() *Ctrl {
	return &Ctrl{}
}

func (c *Ctrl) SetHashiraClient(cli *hashirac.Client) {
	c.hashirac = cli
}

func (c *Ctrl) SetPublisher(p Publisher) {
	c.pub = p
}

func (c *Ctrl) ConfigureKeyBindings(g *gocui.Gui) error {
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.Quit)
}

func (c *Ctrl) Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Ctrl) list(ctx context.Context) error {
	_, err := c.hashirac.Retrieve(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve task list")
	}
}
