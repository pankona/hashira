package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/pankona/hashira/service"
)

// View represents a view of hashira-cui's mvc
type View struct {
	panes        map[string]*Pane
	gui          *gocui.Gui
	cursor       *cursor
	focusedIndex int
	selectedTask *service.Task
	editingTask  *service.Task
	priorities   []*service.Priority
	pane         *Pane // for restoring focused pane after input
	Delegater
}

// pane names
var pn = []string{
	"Backlog",
	"To Do",
	"Doing",
	"Done",
}

// Initialize initializes view
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

	v.gui = g
	v.Delegater = d
	v.cursor = &cursor{
		index:       0,
		focusedPane: v.panes[pn[0]],
	}

	return nil
}

// Left represents action for left key
// TODO: should be more suitable name
func (v *View) Left() error {
	dst := v.panes[v.gui.CurrentView().Name()].left
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

// Right represents action for right key
// TODO: should be more suitable name
func (v *View) Right() error {
	dst := v.panes[v.gui.CurrentView().Name()].right
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

	return v.Delegate(UpdateBulk, t, v.priorities)
}

func (v *View) changeFocusedPane(pane *Pane) error {
	v.cursor.focusedPane = pane

	// Resume scroll status
	// Check if index is positive since there's possibility that
	// this variable becomes negative if pane has no task
	if v.cursor.index >= 0 {
		v.focusedIndex = pane.renderFrom + v.cursor.index
	}

	_, err := v.gui.SetCurrentView(pane.name)
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

// Up represents action for up key
// TODO: should be more suitable name
func (v *View) Up(g *gocui.Gui, _ *gocui.View) error {
	v.cursor.index--
	v.focusedIndex--

	if v.selectedTask == nil {
		return nil
	}
	return v.setPriorityHigh(v.priorities, v.selectedTask)
}

// Down represents action for down key
// TODO: should be more suitable name
func (v *View) Down(g *gocui.Gui, _ *gocui.View) error {
	v.cursor.index++
	v.focusedIndex++

	if v.selectedTask == nil {
		return nil
	}
	return v.setPriorityLow(v.selectedTask)
}

func (v *View) setPriorityHigh(priorities []*service.Priority, task *service.Task) error {
	p, err := v.lookupPaneByTask(task)
	if err != nil {
		return err
	}

	ids := p.priorities

	for i, id := range ids {
		if id == task.Id {
			if i == 0 {
				// already on top. do nothing.
				return nil
			}
			// swap
			ids[i-1], ids[i] = ids[i], ids[i-1]
			return nil
		}
	}

	return fmt.Errorf("failed to set priority high for task [%s]", task.Name)
}

func (v *View) setPriorityLow(task *service.Task) error {
	p, err := v.lookupPaneByTask(task)
	if err != nil {
		return err
	}

	ids := p.priorities

	for i, id := range ids {
		if id == task.Id {
			if i == len(ids)-1 {
				// already on bottom. do nothing.
				return nil
			}
			// swap
			ids[i+1], ids[i] = ids[i], ids[i+1]
			return nil
		}
	}

	return fmt.Errorf("failed to set priority low for task [%s]", task.Name)
}

// markTaskAsDone moves specified task to Done pane.
// If the specified task is already on Done, the task is deleted.
func (v *View) markTaskAsDone(t *service.Task) error {
	p := v.panes[pn[len(pn)-1]] // last pane (may be Done)
	if t == nil || p == nil {
		return nil
	}

	if t == v.selectedTask {
		// deselect
		v.selectedTask = nil
	}

	if t.Place == p.place {
		t.IsDeleted = true
		return v.Delegate(UpdateTask, t)
	}

	return v.moveTaskPlaceTo(t, p, 0)
}

// selectFocusedTask selects focused task.
// call this function again for deselect (toggle).
func (v *View) selectFocusedTask() error {
	if v.selectedTask != nil {
		v.selectedTask = nil
		// on deselect task, it means the deselected task's
		// priority is determined. update priority is necessary.
		return v.Delegate(UpdatePriority, v.priorities)
	}

	v.selectedTask = v.FocusedTask()
	return nil
}

// FocusedTask returns a task that is focused by cursor
func (v *View) FocusedTask() *service.Task {
	if v.focusedIndex < 0 {
		return nil
	}
	id := v.cursor.focusedPane.priorities[v.focusedIndex]
	return v.cursor.focusedPane.tasks[id]
}

type directive int

const (
	dirRight directive = iota
	dirLeft
)

func (v *View) lookupPaneByTask(t *service.Task) (*Pane, error) {
	for i, p := range v.panes {
		if p.place == t.Place {
			return v.panes[i], nil
		}
	}
	return nil, fmt.Errorf("failed to lookup pane by task")
}

// TODO: Performance concern.
// This function executes two communication (over delegation),
// in case that calling this function quickly and continuously, it causes
// awkward movement.
func (v *View) moveTaskTo(t *service.Task, dir directive) error {
	pane, err := v.lookupPaneByTask(t)
	if err != nil {
		return err
	}

	switch dir {
	case dirRight:
		t.Place = pane.right.place
		pane = pane.right
	case dirLeft:
		t.Place = pane.left.place
		pane = pane.left
	}

	err = v.Delegate(UpdateTask, t)
	if err != nil {
		return fmt.Errorf("failed to update: %s", err.Error())
	}

	// put the moved task on top of pane
	// TODO: don't specify 0 as cap to improve performance
	priorities := make([]string, 0)
	priorities = append([]string{t.Id}, priorities...)
	for _, id := range pane.priorities {
		if t.Id != id {
			priorities = append(priorities, id)
		}
	}

	// update priorities by newly created priority array
	pane.priorities = priorities
	for i, p := range v.priorities {
		if p.Place == pane.place {
			v.priorities[i] = &service.Priority{
				Place: pane.place,
				Ids:   pane.priorities,
			}
		}
	}

	return v.Delegate(UpdatePriority, v.priorities)
}

func (v *View) input(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		return v.determineInput(g, gv)
	}
	return v.showInput(g)
}

func (v *View) showInput(g *gocui.Gui) error {
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
		// TODO: should inject Editor interface
		_ = g.SetKeybinding(input.Name(), gocui.KeyCtrlH, gocui.ModNone, v.KeyCtrlH)
		_ = g.SetKeybinding(input.Name(), gocui.KeyCtrlL, gocui.ModNone, v.KeyCtrlL)
		g.Cursor = true
	}

	// use this pane to restore focused pane after input
	v.pane = v.panes[g.CurrentView().Name()]

	_, err = g.SetCurrentView("input")
	return err
}

func (v *View) determineInput(g *gocui.Gui, gv *gocui.View) error {
	defer func() {
		v.editingTask = nil
		g.DeleteKeybindings(gv.Name())
		g.Cursor = false

		err := g.DeleteView(gv.Name())
		if err != nil {
			log.Printf("[WARNING] failed to delete view: %v", err)
		}

		if v.pane == nil {
			// should not reach.
			log.Printf("[WARNING] pane to restore after input is nil")
			v.pane = v.panes[pn[0]]
		}

		_, err = g.SetCurrentView(v.pane.name)
		if err != nil {
			log.Printf("[WARNING] failed to restore current view: %v", err)
		}

		v.pane = nil
	}()

	msg := gv.Buffer()
	msg = strings.TrimSuffix(msg, "\n")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return nil
	}

	if v.editingTask == nil {
		t := &service.Task{
			Name:  msg,
			Place: v.pane.place,
		}
		return v.Delegate(AddTask, t)
	}
	v.editingTask.Name = msg
	return v.Delegate(UpdateTask, v.editingTask)
}

// Quit quits hashira-cui application
// TODO: should be more suitable name
func (v *View) Quit(g *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

var once = sync.Once{}

// Layout renders panes on screen
func (v *View) Layout(g *gocui.Gui) error {
	for _, p := range v.panes {

		focusedIndex := -1
		if p == v.cursor.focusedPane {
			if v.focusedIndex < 0 {
				v.focusedIndex = 0
			}
			if len(v.cursor.focusedPane.priorities)-1 < v.focusedIndex {
				v.focusedIndex = len(v.cursor.focusedPane.priorities) - 1
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

// OnEvent is called when PubSub publisher publishes messages.
// This method is necessary to fulfill Subscriber interface.
func (v *View) OnEvent(event string, data ...interface{}) {
	v.gui.Update(func(*gocui.Gui) error {
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

		return nil
	})
}
