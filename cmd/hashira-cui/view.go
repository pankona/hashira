package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/pankona/gocui"
	"github.com/pankona/hashira/service"
	"github.com/pankona/orderedmap"
)

// View represents a view of hashira-cui's mvc
type View struct {
	panes        map[string]*Pane
	gui          *gocui.Gui
	cursor       *cursor
	focusedIndex int
	selectedTask *KeyedTask
	editingTask  *KeyedTask
	priorities   map[string]*service.Priority
	pane         *Pane     // for restoring focused pane after input
	pn           []string  // pane names
	once         sync.Once // for set current view on first Layout
	editor       gocui.Editor
	Delegater
}

// Initialize initializes view
func (v *View) Initialize(g *gocui.Gui, d Delegater) {
	// pane display names
	v.pn = []string{
		"Backlog",
		"To Do",
		"Doing",
		"Done",
	}

	v.panes = make(map[string]*Pane)

	v.panes[v.pn[0]] = &Pane{
		name:  v.pn[0],
		index: 0,
		place: service.Place_BACKLOG,
	}
	v.panes[v.pn[1]] = &Pane{
		name:  v.pn[1],
		index: 1,
		place: service.Place_TODO,
	}
	v.panes[v.pn[2]] = &Pane{
		name:  v.pn[2],
		index: 2,
		place: service.Place_DOING,
	}
	v.panes[v.pn[3]] = &Pane{
		name:  v.pn[3],
		index: 3,
		place: service.Place_DONE,
	}

	for k := range v.panes {
		v.panes[k].tasks = orderedmap.New()
	}

	v.panes[v.pn[0]].right = v.panes[v.pn[1]]
	v.panes[v.pn[1]].right = v.panes[v.pn[2]]
	v.panes[v.pn[2]].right = v.panes[v.pn[3]]
	v.panes[v.pn[3]].right = v.panes[v.pn[0]]

	v.panes[v.pn[0]].left = v.panes[v.pn[3]]
	v.panes[v.pn[1]].left = v.panes[v.pn[0]]
	v.panes[v.pn[2]].left = v.panes[v.pn[1]]
	v.panes[v.pn[3]].left = v.panes[v.pn[2]]

	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue

	v.gui = g
	v.Delegater = d
	v.cursor = &cursor{
		index:       0,
		focusedPane: v.panes[v.pn[0]],
	}
	v.editor = &hashiraEditor{}
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

func (v *View) moveTaskPlaceTo(t *KeyedTask, pane *Pane, insertTo int) error {
	oldhaunt, err := v.lookupPaneByTask(t)
	if err != nil {
		log.Printf("failed to lookup pane by task: %v", err)
	}

	t.Place = pane.place

	// remove specified task from the task originally belonged to
	err = oldhaunt.tasks.RemoveByKey(t.Id)
	if err != nil {
		log.Printf("failed to remove a task [%s] from [%s]: %v", t.Id, oldhaunt.name, err)
	}

	err = pane.tasks.Insert(t, insertTo)
	if err != nil {
		log.Printf("failed to insert [%s] to [%s]. fatal: %v", t.Name, pane.name, err)
	}

	v.priorities[oldhaunt.place.String()].Ids = oldhaunt.tasks.Order()
	v.priorities[pane.place.String()].Ids = pane.tasks.Order()

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

// Up represents action for up key
// TODO: should be more suitable name
func (v *View) Up(g *gocui.Gui, _ *gocui.View) error {
	v.cursor.index--
	v.focusedIndex--

	if v.selectedTask == nil {
		return nil
	}

	return v.setPriorityHigh(v.selectedTask)
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

func (v *View) setPriorityHigh(task *KeyedTask) error {
	p, err := v.lookupPaneByTask(task)
	if err != nil {
		return err
	}

	ids := p.tasks.Order()

	for i, id := range ids {
		if id == task.Id {
			if i == 0 {
				// already on top. do nothing.
				return nil
			}
			// swap
			p.tasks.Swap(i-1, i)
			return nil
		}
	}

	return fmt.Errorf("failed to set priority high for task [%s]", task.Name)
}

func (v *View) setPriorityLow(task *KeyedTask) error {
	p, err := v.lookupPaneByTask(task)
	if err != nil {
		return err
	}

	ids := p.tasks.Order()

	for i, id := range ids {
		if id == task.Id {
			if i == len(ids)-1 {
				// already on bottom. do nothing.
				return nil
			}
			// swap
			p.tasks.Swap(i, i+1)
			return nil
		}
	}

	return fmt.Errorf("failed to set priority low for task [%s]", task.Name)
}

// markTaskAsDone moves specified task to Done pane.
// If the specified task is already on Done, the task is deleted.
func (v *View) markTaskAsDone(t *KeyedTask) error {
	p := v.panes[v.pn[len(v.pn)-1]] // last pane (may be Done)
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
		// on deselect task, it means the deselected task's
		// priority is determined. update priority is necessary.
		p, err := v.lookupPaneByTask(v.selectedTask)
		if err != nil {
			log.Printf("[WARNING] failed to lookup pane by specified task [%s:%s]",
				v.selectedTask.Id, v.selectedTask.Name)
		}

		v.priorities[p.place.String()] = &service.Priority{
			Ids: p.tasks.Order(),
		}

		v.selectedTask = nil

		return v.Delegate(UpdatePriority, v.priorities)
	}

	v.selectedTask = v.FocusedTask()
	return nil
}

// FocusedTask returns a task that is focused by cursor
func (v *View) FocusedTask() *KeyedTask {
	if v.focusedIndex < 0 ||
		v.cursor.focusedPane.tasks.Len() <= v.focusedIndex {
		return nil
	}
	t := v.cursor.focusedPane.tasks.GetByIndex(v.focusedIndex)
	return t.(*KeyedTask)
}

type direction int

const (
	dirRight direction = iota
	dirLeft
)

func (v *View) lookupPaneByTask(t *KeyedTask) (*Pane, error) {
	for k, p := range v.panes {
		if p.place == t.Place {
			return v.panes[k], nil
		}
	}
	return nil, fmt.Errorf("failed to lookup pane by task")
}

func (v *View) moveTaskTo(t *KeyedTask, dir direction) error {
	pane, err := v.lookupPaneByTask(t)
	if err != nil {
		return err
	}

	switch dir {
	case dirRight:
		pane = pane.right
	case dirLeft:
		pane = pane.left
	}

	return v.moveTaskPlaceTo(t, pane, 0)
}

func (v *View) input(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == inputViewName {
		return v.determineInput(g, gv)
	}
	return v.showInput(g)
}

func (v *View) hideInput(g *gocui.Gui, gv *gocui.View) error {
	v.editingTask = nil
	g.DeleteKeybindings(gv.Name())
	g.Cursor = false

	err := g.DeleteView(gv.Name())
	if err != nil {
		return fmt.Errorf("[WARNING] failed to delete view: %v", err)
	}

	if v.pane == nil {
		// should not reach.
		log.Printf("[WARNING] pane to restore after input is nil")
		v.pane = v.panes[v.pn[0]]
	}

	_, err = g.SetCurrentView(v.pane.name)
	if err != nil {
		return fmt.Errorf("[WARNING] failed to restore current view: %v", err)
	}

	v.pane = nil
	return nil
}

func (v *View) showInput(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	input, err := g.SetView(inputViewName, maxX/2-20, maxY/2, maxX/2+20, maxY/2+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if v.editingTask == nil {
			input.Title = "New task?"
		} else {
			input.Title = "Update task?"
			_, err = fmt.Fprintf(input, "%s", v.editingTask.Name)
			if err != nil {
				return fmt.Errorf("failed to write on input view for update: %s", err)
			}
		}
		input.Editable = true
		input.Editor = v.editor
		input.MoveCursor(runewidth.StringWidth(input.Buffer()), 0, true)
		g.Cursor = true
	}

	// use this pane to restore focused pane after input
	v.pane = v.panes[g.CurrentView().Name()]

	_, err = g.SetCurrentView(inputViewName)
	return err
}

func (v *View) determineInput(g *gocui.Gui, gv *gocui.View) error {
	defer func() {
		err := v.hideInput(g, gv)
		if err != nil {
			log.Printf("[WARNING] failed to hide input: %v", err)
		}
	}()

	msg := gv.Buffer()
	msg = strings.TrimSuffix(msg, "\n")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return nil
	}

	if v.editingTask == nil {
		t := &KeyedTask{
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

// Layout renders panes on screen
func (v *View) Layout(g *gocui.Gui) error {
	for _, p := range v.panes {
		focusedIndex := -1
		if p == v.cursor.focusedPane {
			if v.focusedIndex < 0 {
				v.focusedIndex = 0
			}
			if v.cursor.focusedPane.tasks.Len()-1 < v.focusedIndex {
				v.focusedIndex = v.cursor.focusedPane.tasks.Len() - 1
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
	v.once.Do(func() {
		if _, err := g.SetCurrentView(v.pn[0]); err != nil {
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
			tasks := data[0].(map[string]*KeyedTask)
			v.priorities = data[1].(map[string]*service.Priority)

			for i := range v.panes {
				priority := v.priorities[v.panes[i].place.String()].Ids

				// reset tasks
				v.panes[i].tasks = orderedmap.New()

				for _, id := range priority {
					t := tasks[id]
					err := v.panes[i].tasks.Add(t)
					if err != nil {
						log.Printf("failed to add a task [%s:%s]. skip: %v", t.Id, t.Name, err)
					}
				}
			}
		default:
			// nop
		}

		return nil
	})
}
