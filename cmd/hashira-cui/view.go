package main

import (
	"fmt"

	"github.com/pankona/hashira/service"

	"github.com/jroimartin/gocui"
)

type View struct {
	pains map[string]*Pane
	g     *gocui.Gui
}

type Pane struct {
	name  string
	index int // place of this pane
	left  *Pane
	right *Pane
	tasks []*service.Task
}

func (v *View) Initialize(g *gocui.Gui) error {
	v.pains = make(map[string]*Pane)

	v.pains["Backlog"] = &Pane{name: "Backlog", index: 0}
	v.pains["To Do"] = &Pane{name: "To Do", index: 1}
	v.pains["Doing"] = &Pane{name: "Doing", index: 2}
	v.pains["Done"] = &Pane{name: "Done", index: 3}

	v.pains["Backlog"].right = v.pains["To Do"]
	v.pains["To Do"].right = v.pains["Doing"]
	v.pains["Doing"].right = v.pains["Done"]
	v.pains["Done"].right = v.pains["Backlog"]

	v.pains["Backlog"].left = v.pains["Done"]
	v.pains["To Do"].left = v.pains["Doing"]
	v.pains["Doing"].left = v.pains["To Do"]
	v.pains["Done"].left = v.pains["Backlog"]

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue
	g.SetCurrentView(v.pains["Backlog"].name)

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
	for _, v := range v.pains {
		err := v.Layout(g)
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

var place = map[int]service.Place{
	0: service.Place_BACKLOG,
	1: service.Place_TODO,
	2: service.Place_DOING,
	3: service.Place_DONE,
}

func (p *Pane) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(p.name, maxX/4*p.index, 1, maxX/4*p.index+maxX/4-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = p.name
		for _, task := range p.tasks {
			if task.Place == place[p.index] && !task.IsDeleted {
				_, err = fmt.Fprintln(v, task.Name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
