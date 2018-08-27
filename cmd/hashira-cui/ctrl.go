package main

import (
	"context"

	"github.com/pankona/hashira/service"
	"github.com/pkg/errors"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

type Ctrl struct {
	m   *Model
	pub Publisher
}

func NewCtrl(m *Model) *Ctrl {
	return &Ctrl{m: m}
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

func (c *Ctrl) Update(ctx context.Context) error {
	tasks, err := c.m.List(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to update tasks")
	}

	c.pub.Publish("update", tasks)
}

type Model struct {
	hashirac *hashirac.Client
}

func (m *Model) SetHashiraClient(cli *hashirac.Client) {
	m.hashirac = cli
}

func (m *Model) List(ctx context.Context) ([]*service.Task, error) {
	return m.hashirac.Retrieve(ctx)
}
