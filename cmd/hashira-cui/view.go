package main

import (
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

type View struct {
	pains         map[string]*Pane
	g             *gocui.Gui
	selectedIndex int
}

// pane names
var pn = []string{
	"Backlog",
	"To Do",
	"Doing",
	"Done",
}

func (v *View) Initialize(g *gocui.Gui) error {
	v.pains = make(map[string]*Pane)

	v.pains[pn[0]] = &Pane{
		name:  pn[0],
		index: 0, place: service.Place_BACKLOG}
	v.pains[pn[1]] = &Pane{
		name:  pn[1],
		index: 1, place: service.Place_TODO}
	v.pains[pn[2]] = &Pane{
		name:  pn[2],
		index: 2, place: service.Place_DOING}
	v.pains[pn[3]] = &Pane{
		name:  pn[3],
		index: 3, place: service.Place_DONE}

	v.pains[pn[0]].right = v.pains[pn[1]]
	v.pains[pn[1]].right = v.pains[pn[2]]
	v.pains[pn[2]].right = v.pains[pn[3]]
	v.pains[pn[3]].right = v.pains[pn[0]]

	v.pains[pn[0]].left = v.pains[pn[3]]
	v.pains[pn[1]].left = v.pains[pn[0]]
	v.pains[pn[2]].left = v.pains[pn[1]]
	v.pains[pn[3]].left = v.pains[pn[2]]

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue
	g.SetCurrentView(v.pains[pn[0]].name)

	v.g = g

	return nil
}

func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	for _, p := range v.pains {
		_ = g.SetKeybinding(p.name, 'h', gocui.ModNone, v.Left)
		_ = g.SetKeybinding(p.name, 'l', gocui.ModNone, v.Right)
		_ = g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up)
		_ = g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down)
	}
	_ = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, v.Enter)
	_ = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit)
	return nil
}

func (v *View) Left(g *gocui.Gui, _ *gocui.View) error {
	g.SetCurrentView(v.pains[g.CurrentView().Name()].left.name)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

func (v *View) Right(g *gocui.Gui, _ *gocui.View) error {
	g.SetCurrentView(v.pains[g.CurrentView().Name()].right.name)
	g.Update(func(*gocui.Gui) error {
		return nil
	})
	return nil
}

func (v *View) Up(g *gocui.Gui, _ *gocui.View) error {
	v.selectedIndex--
	return nil
}

func (v *View) Down(g *gocui.Gui, _ *gocui.View) error {
	v.selectedIndex++
	return nil
}

func (v *View) Enter(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		err := g.DeleteView("input")
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView(v.pains[pn[0]].name)
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("input", maxX/2-20, maxY/2, maxX/2+20, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "New task?"
		v.Editable = true
		g.SetCurrentView("input")
	}
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

var once = sync.Once{}

func (v *View) Layout(g *gocui.Gui) error {
	for _, p := range v.pains {
		if g.CurrentView() != nil &&
			g.CurrentView().Name() == p.name {
			if v.selectedIndex <= 0 {
				v.selectedIndex = 0
			}
			if v.selectedIndex >= p.len() {
				v.selectedIndex = p.len() - 1
			}
		}

		err := p.Layout(g, v.selectedIndex)
		if err != nil {
			return err
		}
	}

	// initialize current view
	// this function only needs to be called once on starting application
	once.Do(func() {
		if _, err := g.SetCurrentView(pn[0]); err != nil {
			panic(err)
		}
	})

	return nil
}

func (v *View) OnEvent(event string, data interface{}) {
	switch event {
	case "update":
		tasks := data.([]*service.Task)
		for i := range v.pains {
			// TODO: filter tasks for each pane
			v.pains[i].tasks = tasks
		}

		// ZZZ: this update may not be needed
		v.g.Update(func(*gocui.Gui) error {
			return nil
		})
	default:
		// nop
	}
}
