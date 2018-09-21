package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

type View struct {
	pains         map[string]*Pane
	g             *gocui.Gui
	selectedIndex int
	grabbedTask   *service.Task
	editingTask   *service.Task
	priorities    []*service.Priority
	Delegater
}

type Delegater interface {
	Delegate(event string, data interface{}) error
}

// pane names
var pn = []string{
	"Backlog",
	"To Do",
	"Doing",
	"Done",
}

func (v *View) Initialize(g *gocui.Gui, d Delegater) error {
	v.pains = make(map[string]*Pane)

	v.pains[pn[0]] = &Pane{
		name:  pn[0],
		index: 0,
		place: service.Place_BACKLOG,
	}
	v.pains[pn[1]] = &Pane{
		name:  pn[1],
		index: 1,
		place: service.Place_TODO,
	}
	v.pains[pn[2]] = &Pane{
		name:  pn[2],
		index: 2, place: service.Place_DOING,
	}
	v.pains[pn[3]] = &Pane{
		name:  pn[3],
		index: 3, place: service.Place_DONE,
	}

	for k := range v.pains {
		v.pains[k].tasks = make(map[string]*service.Task)
	}

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
	v.Delegater = d

	return nil
}

func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	for _, p := range v.pains {
		_ = g.SetKeybinding(p.name, 'h', gocui.ModNone, v.Left)   // TODO: should be v.KeyH
		_ = g.SetKeybinding(p.name, 'l', gocui.ModNone, v.Right)  // TODO: should be v.KeyL
		_ = g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up)     // TODO: should be v.KeyK
		_ = g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down)   // TODO: should be v.KeyJ
		_ = g.SetKeybinding(p.name, 'x', gocui.ModNone, v.Delete) // TODO: should be v.KeyX
		_ = g.SetKeybinding(p.name, gocui.KeySpace, gocui.ModNone, v.KeySpace)
		_ = g.SetKeybinding(p.name, 'i', gocui.ModNone, v.KeyI)
		_ = g.SetKeybinding(p.name, 'e', gocui.ModNone, v.KeyE)
	}
	_ = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, v.KeyEnter)
	_ = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit) // TODO: should be v.KeyCtrlC
	return nil
}

func (v *View) KeyE(g *gocui.Gui, gv *gocui.View) error {
	t := v.SelectedItem()
	if t == nil {
		return nil
	}

	v.editingTask = t

	return v.input(g, gv)
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
	if v.grabbedTask != nil {
		v.setPriorityHigh(v.priorities, v.grabbedTask)
	}
	v.selectedIndex--
	return nil
}

func (v *View) Down(g *gocui.Gui, _ *gocui.View) error {
	if v.grabbedTask != nil {
		v.setPriorityLow(v.priorities, v.grabbedTask)
	}
	v.selectedIndex++
	return nil
}

// TODO: should be delegated
func (v *View) setPriorityHigh(priorities []*service.Priority, task *service.Task) error {
	place := task.Place
	var index int
	for i, p := range v.priorities {
		if p.Place == place {
			index = i
		}
	}

	log.Printf("priority updated (before): %v", v.priorities[index])

	for i, id := range v.priorities[index].Ids {
		if id == task.Id {
			if i == 0 {
				return nil
			}
			v.priorities[index].Ids[i-1], v.priorities[index].Ids[i] = v.priorities[index].Ids[i], v.priorities[index].Ids[i-1]
			break
		}
	}

	log.Printf("priority updated (after): %v", v.priorities[index])

	return nil
}

// TODO: should be delegated
func (v *View) setPriorityLow(priorities []*service.Priority, task *service.Task) error {
	place := task.Place
	var index int
	for i, p := range v.priorities {
		if p.Place == place {
			index = i
		}
	}

	for i, id := range v.priorities[index].Ids {
		if id == task.Id {
			if i == len(v.priorities[index].Ids)-1 {
				return nil
			}
			v.priorities[index].Ids[i+1], v.priorities[index].Ids[i] = v.priorities[index].Ids[i], v.priorities[index].Ids[i+1]
			break
		}
	}

	return nil
}

func (v *View) Delete(g *gocui.Gui, _ *gocui.View) error {
	// TODO: (word revision) "Selected" should be "Focused"
	t := v.SelectedItem()
	if t == nil {
		return nil
	}
	return v.Delegate("delete", t)
}

func (v *View) KeySpace(g *gocui.Gui, gv *gocui.View) error {
	return v.grabFocusedItem()
}

func (v *View) grabFocusedItem() error {
	if v.grabbedTask != nil {
		v.grabbedTask = nil
		v.Delegate("updatePriority", v.priorities)
	} else {
		v.grabbedTask = v.SelectedItem()
	}
	return nil
}

func (v *View) SelectedItem() *service.Task {
	if v.selectedIndex < 0 {
		return nil
	}
	p := v.pains[v.g.CurrentView().Name()]
	// FIXME
	log.Printf("selectedIndex = %d, priorities: %v", v.selectedIndex, v.priorities[0].Ids)
	var index int
	for i, t := range v.priorities {
		if t.Place == p.place {
			index = i
		}
	}
	id := v.priorities[index].Ids[v.selectedIndex]
	return p.tasks[id]
}

func (v *View) KeyI(g *gocui.Gui, gv *gocui.View) error {
	return v.input(g, gv)
}

func (v *View) KeyEnter(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		return v.input(g, gv)
	}

	// TODO: implement
	return nil
}

func (v *View) input(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		defer func() {
			v.editingTask = nil
			g.Cursor = false
			g.DeleteView("input")
			// TODO: set selected pane as current view
			g.SetCurrentView(v.pains[pn[0]].name)
		}()

		msg := gv.Buffer()
		if msg == "" {
			return nil
		}
		msg = strings.TrimSuffix(msg, "\n")

		if v.editingTask == nil {
			return v.Delegate("add", msg)
		}
		v.editingTask.Name = msg
		return v.Delegate("update", v.editingTask)
	}

	maxX, maxY := g.Size()
	if input, err := g.SetView("input", maxX/2-20, maxY/2, maxX/2+20, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if v.editingTask == nil {
			input.Title = "New task?"
		} else {
			input.Title = "Update task?"
			_, err = fmt.Fprintf(input, v.editingTask.Name)
			if err != nil {
				return fmt.Errorf("failed to write on input view for update: %s", err)
			}
		}
		input.Editable = true
		input.MoveCursor(len(input.Buffer())-1, 0, true)
		g.Cursor = true
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
			log.Printf("p.len() = %d", p.len())
			if v.selectedIndex >= p.len() {
				v.selectedIndex = p.len() - 1
			}
		}

		err := p.Layout(g, v.selectedIndex, v.grabbedTask)
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

func (v *View) OnEvent(event string, data ...interface{}) {
	switch event {
	case "update":
		tasks := data[0].([]*service.Task)
		v.priorities = data[1].([]*service.Priority)
		for _, p := range v.priorities {
			log.Printf("view receives priority of : %s", p.Place.String())
			log.Printf("view receives priority : %v", p)
		}

		for i := range v.pains {
			// reset tasks
			v.pains[i].tasks = make(map[string]*service.Task)
			for _, t := range tasks {
				if v.pains[i].place == t.Place {
					v.pains[i].tasks[t.Id] = t
				}
			}

			for _, p := range v.priorities {
				if v.pains[i].place == p.Place {
					v.pains[i].priorities = p.Ids
				}
			}
		}
	default:
		// nop
	}
}
