package main

import (
	"fmt"

	"github.com/pankona/hashira/service"

	"github.com/jroimartin/gocui"
)

type View struct {
	pains [4]Pain
	g     *gocui.Gui
}

func (v *View) Initialize(g *gocui.Gui) error {
	v.pains[0].name = "Backlog"
	v.pains[1].name = "To Do"
	v.pains[2].name = "Doing"
	v.pains[3].name = "Done"

	v.pains[0].right = &v.pains[1]
	v.pains[1].right = &v.pains[2]
	v.pains[2].right = &v.pains[3]
	v.pains[3].right = &v.pains[0]

	v.pains[0].left = &v.pains[3]
	v.pains[1].left = &v.pains[2]
	v.pains[2].left = &v.pains[1]
	v.pains[3].left = &v.pains[0]

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue
	g.SetCurrentView(v.pains[0].name)

	v.g = g

	return nil
}

func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	_ = g.SetKeybinding("Backlog", 'h', gocui.ModNone, v.Left)
	_ = g.SetKeybinding("Backlog", 'l', gocui.ModNone, v.Right)
	_ = g.SetKeybinding("To Do", 'h', gocui.ModNone, v.Left)
	_ = g.SetKeybinding("To Do", 'l', gocui.ModNone, v.Right)
	_ = g.SetKeybinding("Doing", 'h', gocui.ModNone, v.Left)
	_ = g.SetKeybinding("Doing", 'l', gocui.ModNone, v.Right)
	_ = g.SetKeybinding("Done", 'h', gocui.ModNone, v.Left)
	_ = g.SetKeybinding("Done", 'l', gocui.ModNone, v.Right)
	_ = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit)
	return nil
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

func (v *View) Left(g *gocui.Gui, _ *gocui.View) error {
	g.SetCurrentView(p[g.CurrentView().Name()].left)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

func (v *View) Right(g *gocui.Gui, _ *gocui.View) error {
	g.SetCurrentView(p[g.CurrentView().Name()].right)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

func (v *View) SetFocus(name string) error {
	_, err := v.g.SetCurrentView(name)
	v.g.Update(func(*gocui.Gui) error { return nil })
	return err
}

func (v *View) Quit(g *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (v *View) Layout(g *gocui.Gui) error {
	for i := range v.pains {
		err := v.pains[i].Layout(i, g)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *View) OnEvent(event string, data interface{}) {
	switch event {
	case "update":
		tasks := data.([]*service.Task)
		for i := range v.pains {
			v.pains[i].tasks = tasks
		}
	default:
		// nop
	}
}

type Pain struct {
	name  string
	left  *Pain
	right *Pain
	tasks []*service.Task
}

var place = map[int]service.Place{
	0: service.Place_BACKLOG,
	1: service.Place_TODO,
	2: service.Place_DOING,
	3: service.Place_DONE,
}

func (p *Pain) Layout(index int, g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(p.name, maxX/4*index, 1, maxX/4*index+maxX/4-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
		for _, task := range p.tasks {
			if task.Place == place[index] && !task.IsDeleted {
				_, err = fmt.Fprintln(v, task.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
