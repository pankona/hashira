package main

import (
	"context"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

type Ctrl struct {
	m *Model
	g *gocui.Gui
}

func (c *Ctrl) Initialize() {
	// TODO: should place on view
	c.g.Highlight = true
	c.g.SelFgColor = gocui.ColorBlue
}

func (c *Ctrl) ConfigureKeyBindings(g *gocui.Gui) error {
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.Quit)
}

func (c *Ctrl) Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Ctrl) SetFocus(name string) error {
	// TODO: should place on view
	_, err := c.g.SetCurrentView(name)
	c.g.Update(func(*gocui.Gui) error { return nil })
	return err
}

func (c *Ctrl) Update(ctx context.Context) error {
	err := c.m.List(ctx)
	if err != nil {
		return err
	}

	// TODO: this is refresh of view. should place on view
	c.g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

type Model struct {
	hashirac *hashirac.Client
	pub      Publisher
}

func (m *Model) SetPublisher(p Publisher) {
	m.pub = p
}

func (m *Model) SetHashiraClient(cli *hashirac.Client) {
	m.hashirac = cli
}

func (m *Model) List(ctx context.Context) error {
	tasks, err := m.hashirac.Retrieve(ctx)
	if err != nil {
		return err
	}
	m.pub.Publish("update", tasks)
	return nil
}
