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
	panes        map[string]*Pane
	g            *gocui.Gui
	cursor       *cursor
	focusedIndex int
	selectedTask *service.Task
	editingTask  *service.Task
	priorities   []*service.Priority
	pane         *Pane // for restoring focused pane after input
	Delegater
}

type cursor struct {
	index int
	pane  *Pane
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
	v.panes = make(map[string]*Pane)

	v.panes[pn[0]] = &Pane{
		name:  pn[0],
		index: 0,
		place: service.Place_BACKLOG,
	}
	v.panes[pn[1]] = &Pane{
		name:  pn[1],
		index: 1,
		place: service.Place_TODO,
	}
	v.panes[pn[2]] = &Pane{
		name:  pn[2],
		index: 2, place: service.Place_DOING,
	}
	v.panes[pn[3]] = &Pane{
		name:  pn[3],
		index: 3, place: service.Place_DONE,
	}

	for k := range v.panes {
		v.panes[k].tasks = make(map[string]*service.Task)
	}

	v.panes[pn[0]].right = v.panes[pn[1]]
	v.panes[pn[1]].right = v.panes[pn[2]]
	v.panes[pn[2]].right = v.panes[pn[3]]
	v.panes[pn[3]].right = v.panes[pn[0]]

	v.panes[pn[0]].left = v.panes[pn[3]]
	v.panes[pn[1]].left = v.panes[pn[0]]
	v.panes[pn[2]].left = v.panes[pn[1]]
	v.panes[pn[3]].left = v.panes[pn[2]]

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue
	g.SetCurrentView(v.panes[pn[0]].name)

	v.g = g
	v.Delegater = d
	v.cursor = &cursor{
		index: 0,
		pane:  v.panes[pn[0]],
	}

	return nil
}

func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	for _, p := range v.panes {
		_ = g.SetKeybinding(p.name, 'h', gocui.ModNone, v.Left)  // TODO: should be v.KeyH
		_ = g.SetKeybinding(p.name, 'l', gocui.ModNone, v.Right) // TODO: should be v.KeyL
		_ = g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up)    // TODO: should be v.KeyK
		_ = g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down)  // TODO: should be v.KeyJ
		_ = g.SetKeybinding(p.name, 'x', gocui.ModNone, v.KeyX)
		_ = g.SetKeybinding(p.name, 'i', gocui.ModNone, v.KeyI)
		_ = g.SetKeybinding(p.name, 'I', gocui.ModNone, v.KeyShiftI)
		_ = g.SetKeybinding(p.name, 'e', gocui.ModNone, v.KeyE)
		_ = g.SetKeybinding(p.name, gocui.KeySpace, gocui.ModNone, v.KeySpace)
	}
	_ = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, v.KeyEnter)
	_ = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit) // TODO: should be v.KeyCtrlC
	return nil
}

func (v *View) KeyE(g *gocui.Gui, gv *gocui.View) error {
	t := v.FocusedTask()
	if t == nil {
		return nil
	}

	v.editingTask = t

	return v.input(g, gv)
}

func (v *View) Left(g *gocui.Gui, _ *gocui.View) error {
	dst := v.panes[g.CurrentView().Name()].left
	err := v.changeFocusedPane(dst)
	if err != nil {
		return err
	}

	t := v.selectedTask
	if t != nil {
		return v.moveTaskPlaceTo(t, dst, v.cursor.index)
	}
	return nil
}

func (v *View) Right(g *gocui.Gui, _ *gocui.View) error {
	dst := v.panes[g.CurrentView().Name()].right
	err := v.changeFocusedPane(dst)
	if err != nil {
		return err
	}

	t := v.selectedTask
	if t != nil {
		return v.moveTaskPlaceTo(t, dst, v.cursor.index)
	}
	return nil
}

func (v *View) moveTaskPlaceTo(t *service.Task, pane *Pane, insertTo int) error {
	t.Place = pane.place
	err := v.Delegate("update", t)
	if err != nil {
		return err
	}

	priority := remove(pane.priorities, t.Id)
	priority = insert(priority, t.Id, insertTo)
	if priority == nil {
		return fmt.Errorf("failed to insert [%s] to [%s]. fatal", t.Name, pane.name)
	}

	for i, p := range v.priorities {
		if p.Place == pane.place {
			v.priorities[i].Ids = priority
		}
	}

	return v.Delegate("updatePriority", v.priorities)
}

func (v *View) changeFocusedPane(pane *Pane) error {
	v.cursor.pane = pane
	if v.cursor.index > len(v.cursor.pane.priorities)-1 {
		v.cursor.index = len(v.cursor.pane.priorities) - 1
	}

	// resume scroll status
	var index int
	if v.cursor.index > 0 {
		index = v.cursor.index
	}
	v.focusedIndex = pane.renderFrom + index

	_, err := v.g.SetCurrentView(pane.name)
	return err
}

func remove(ss []string, s string) []string {
	ret := make([]string, len(ss))
	var index int
	var found bool

	for i := 0; i < len(ss); i++ {
		if ss[i] == s {
			found = true
			continue
		}
		ret[index] = ss[i]
		index++
	}

	if found {
		return ret[:len(ss)-1]
	}
	return ret
}

func insert(ss []string, s string, index int) []string {
	if index < 0 {
		return nil
	}
	if index > len(ss) {
		index = len(ss)
	}

	ret := make([]string, len(ss)+1)
	for i := 0; i < index; i++ {
		ret[i] = ss[i]
	}
	ret[index] = s
	for i := index; i < len(ss); i++ {
		ret[i+1] = ss[i]
	}
	return ret
}

func (v *View) Up(g *gocui.Gui, _ *gocui.View) error {
	if v.selectedTask != nil {
		v.setPriorityHigh(v.priorities, v.selectedTask)
	}

	v.cursor.index--
	v.focusedIndex--

	return nil
}

func (v *View) Down(g *gocui.Gui, _ *gocui.View) error {
	if v.selectedTask != nil {
		v.setPriorityLow(v.priorities, v.selectedTask)
	}

	v.cursor.index++
	v.focusedIndex++

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

	for i, id := range v.priorities[index].Ids {
		if id == task.Id {
			if i == 0 {
				return nil
			}
			v.priorities[index].Ids[i-1], v.priorities[index].Ids[i] = v.priorities[index].Ids[i], v.priorities[index].Ids[i-1]
			break
		}
	}

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

func (v *View) KeyX(*gocui.Gui, *gocui.View) error {
	t := v.FocusedTask()
	return v.markTaskAsDone(t)
}

// markTaskAsDone moves specified task to Done pane.
// If the specified task is already on Done, the task is deleted.
func (v *View) markTaskAsDone(t *service.Task) error {
	p := v.panes[pn[len(pn)-1]] // last pane (may be Done)
	if t == nil || p == nil {
		return nil
	}

	if t.Place == p.place {
		t.IsDeleted = true
		return v.Delegate("update", t)
	}

	return v.moveTaskPlaceTo(t, p, 0)
}

func (v *View) KeySpace(g *gocui.Gui, gv *gocui.View) error {
	return v.selectFocusedTask()
}

// selectFocusedTask selects focused task.
// call this function again for deselect.
func (v *View) selectFocusedTask() error {
	if v.selectedTask != nil {
		v.selectedTask = nil
		// on deselect task, it means the deselected task's
		// priority is determined. update priority is necessary.
		return v.Delegate("updatePriority", v.priorities)
	} else {
		v.selectedTask = v.FocusedTask()
	}
	return nil
}

func (v *View) FocusedTask() *service.Task {
	if v.focusedIndex < 0 {
		return nil
	}
	id := v.cursor.pane.priorities[v.focusedIndex]
	return v.cursor.pane.tasks[id]
}

func (v *View) KeyEnter(g *gocui.Gui, gv *gocui.View) error {
	return v.input(g, gv)
}

type directive int

const (
	dirRight directive = iota
	dirLeft
)

func (v *View) KeyI(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		return v.input(g, gv)
	}

	t := v.FocusedTask()
	if t == nil {
		log.Printf("selectedTask is nil")
		return nil
	}
	return v.moveTaskTo(t, dirRight)
}

func (v *View) KeyShiftI(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		return v.input(g, gv)
	}

	t := v.FocusedTask()
	if t == nil {
		log.Printf("selectedTask is nil")
		return nil
	}
	return v.moveTaskTo(t, dirLeft)
}

func (v *View) lookupPaneByTask(t *service.Task) *Pane {
	for i, p := range v.panes {
		if p.place == t.Place {
			return v.panes[i]
		}
	}
	return nil
}

func (v *View) moveTaskTo(t *service.Task, dir directive) error {
	pane := v.lookupPaneByTask(t)
	if pane == nil {
		return fmt.Errorf("couldn't lookup a pane by specified task")
	}

	switch dir {
	case dirRight:
		t.Place = pane.right.place
		pane = pane.right
	case dirLeft:
		t.Place = pane.left.place
		pane = pane.left
	}

	err := v.Delegate("update", t)
	if err != nil {
		return fmt.Errorf("failed to update: %s", err.Error())
	}

	// put the moved task on top of pane
	priorities := make([]string, 0)
	priorities = append([]string{t.Id}, priorities...)
	for _, id := range pane.priorities {
		if t.Id != id {
			priorities = append(priorities, id)
		}
	}

	pane.priorities = priorities
	for i, p := range v.priorities {
		if p.Place == pane.place {
			v.priorities[i] = &service.Priority{
				Place: pane.place,
				Ids:   pane.priorities,
			}
		}
	}

	return v.Delegate("updatePriority", v.priorities)
}

func (v *View) input(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		defer func() {
			v.editingTask = nil
			g.DeleteKeybindings(gv.Name())
			g.Cursor = false
			g.DeleteView(gv.Name())
			if v.pane == nil {
				// should not reach.
				log.Printf("[WARNING] pane to restore after input is nil")
				v.pane = v.panes[pn[0]]
			}
			g.SetCurrentView(v.pane.name)
			v.pane = nil
		}()

		msg := gv.Buffer()
		msg = strings.TrimSuffix(msg, "\n")
		msg = strings.TrimSpace(msg)
		if msg == "" {
			return nil
		}

		if v.editingTask == nil {
			return v.Delegate("add", msg)
		}
		v.editingTask.Name = msg
		return v.Delegate("update", v.editingTask)
	}

	maxX, maxY := g.Size()
	input, err := g.SetView("input", maxX/2-20, maxY/2, maxX/2+20, maxY/2+2)
	if err != nil {
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
		_ = g.SetKeybinding(input.Name(), gocui.KeyCtrlH, gocui.ModNone, v.KeyCtrlH)
		_ = g.SetKeybinding(input.Name(), gocui.KeyCtrlL, gocui.ModNone, v.KeyCtrlL)
		g.Cursor = true
	}

	// use this pane to restore focused pane after input
	v.pane = v.panes[g.CurrentView().Name()]

	_, err = g.SetCurrentView("input")
	return err
}

func (v *View) KeyCtrlH(g *gocui.Gui, gv *gocui.View) error {
	gv.MoveCursor(-1, 0, true)
	return nil
}

func (v *View) KeyCtrlL(g *gocui.Gui, gv *gocui.View) error {
	x, _ := gv.Cursor()
	if len(gv.Buffer())-1 <= x {
		return nil
	}
	gv.MoveCursor(+1, 0, true)
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
	for _, p := range v.panes {

		focusedIndex := -1
		if p == v.cursor.pane {
			if v.focusedIndex < 0 {
				v.focusedIndex = 0
			}
			if len(v.cursor.pane.priorities)-1 < v.focusedIndex {
				v.focusedIndex = len(v.cursor.pane.priorities) - 1
			}
			focusedIndex = v.focusedIndex
		}

		err := p.Layout(g, v.cursor, focusedIndex, v.selectedTask)
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

		for i := range v.panes {
			// reset tasks
			v.panes[i].tasks = make(map[string]*service.Task)
			for _, t := range tasks {
				if v.panes[i].place == t.Place {
					v.panes[i].tasks[t.Id] = t
				}
			}

			for _, p := range v.priorities {
				if v.panes[i].place == p.Place {
					v.panes[i].priorities = p.Ids
				}
			}
		}
	default:
		// nop
	}
}
