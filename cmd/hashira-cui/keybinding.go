package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

// ConfigureKeyBindings configures keybindings for ordinal use
func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	for _, p := range v.panes {
		_ = g.SetKeybinding(p.name, 'h', gocui.ModNone, v.KeyH)
		_ = g.SetKeybinding(p.name, 'l', gocui.ModNone, v.KeyL)
		_ = g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up)   // TODO: should be v.KeyK
		_ = g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down) // TODO: should be v.KeyJ
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

// KeyH reacts for "h"
func (v *View) KeyH(g *gocui.Gui, _ *gocui.View) error {
	return v.Left()
}

// KeyL reacts for "l"
func (v *View) KeyL(g *gocui.Gui, _ *gocui.View) error {
	return v.Right()
}

// KeyX reacts for "x"
func (v *View) KeyX(*gocui.Gui, *gocui.View) error {
	t := v.FocusedTask()
	return v.markTaskAsDone(t)
}

// KeyI reacts for "i"
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

// KeyShiftI reacts for "I"
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

// KeyE reacts for "e"
func (v *View) KeyE(g *gocui.Gui, gv *gocui.View) error {
	t := v.FocusedTask()
	if t == nil {
		return nil
	}

	v.editingTask = t

	return v.input(g, gv)
}

// KeySpace reacts for "Space"
func (v *View) KeySpace(g *gocui.Gui, gv *gocui.View) error {
	return v.selectFocusedTask()
}

// KeyEnter reacts for "Enter"
func (v *View) KeyEnter(g *gocui.Gui, gv *gocui.View) error {
	return v.input(g, gv)
}

// KeyCtrlH reacts when Ctrl-h is pressed while inputting task
// TODO:
// keybindings for input should inject Editor interface
func (v *View) KeyCtrlH(g *gocui.Gui, gv *gocui.View) error {
	gv.MoveCursor(-1, 0, true)
	return nil
}

// KeyCtrlL reacts when Ctrl-l is pressed while inputting task
// TODO:
// keybindings for input should inject Editor interface
func (v *View) KeyCtrlL(g *gocui.Gui, gv *gocui.View) error {
	x, _ := gv.Cursor()
	if len(gv.Buffer())-1 <= x {
		return nil
	}
	gv.MoveCursor(+1, 0, true)
	return nil
}
