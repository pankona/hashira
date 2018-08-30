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

// TODO: should place on view
func (c *Ctrl) ConfigureKeyBindings(g *gocui.Gui) error {
	_ = g.SetKeybinding("Backlog", 'h', gocui.ModNone, c.Left)
	_ = g.SetKeybinding("Backlog", 'l', gocui.ModNone, c.Right)
	_ = g.SetKeybinding("To Do", 'h', gocui.ModNone, c.Left)
	_ = g.SetKeybinding("To Do", 'l', gocui.ModNone, c.Right)
	_ = g.SetKeybinding("Doing", 'h', gocui.ModNone, c.Left)
	_ = g.SetKeybinding("Doing", 'l', gocui.ModNone, c.Right)
	_ = g.SetKeybinding("Done", 'h', gocui.ModNone, c.Left)
	_ = g.SetKeybinding("Done", 'l', gocui.ModNone, c.Right)
	return g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.Quit)
}

var p = map[string]struct {
	left  string
	right string
}{
	"Backlog": {left: "Done", right: "To Do"},
	"To Do":   {left: "Backlog", right: "Doing"},
	"Doing":   {left: "To Do", right: "Done"},
	"Done":    {left: "Doing", right: "Backlog"},
}

// TODO: should place on view
func (c *Ctrl) Left(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView(p[g.CurrentView().Name()].left)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

// TODO: should place on view
func (c *Ctrl) Right(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView(p[g.CurrentView().Name()].right)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

func (c *Ctrl) Quit(*gocui.Gui, *gocui.View) error {
	return gocui.ErrQuit
}

// TODO: should place on view
func (c *Ctrl) SetFocus(name string) error {
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
