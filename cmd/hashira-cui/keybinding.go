package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/jroimartin/gocui"
)

type errContainer struct {
	err error
}

func (e *errContainer) setError(err error) {
	if err != nil {
		e.err = errors.Wrapf(e.err, "failed to set keybinding: %v", err.Error())
	}
}

// ConfigureKeyBindings configures keybindings for ordinal use
func (v *View) ConfigureKeyBindings(g *gocui.Gui) error {
	ec := &errContainer{}
	for _, p := range v.panes {
		ec.setError(g.SetKeybinding(p.name, 'h', gocui.ModNone, v.KeyH))
		ec.setError(g.SetKeybinding(p.name, 'l', gocui.ModNone, v.KeyL))
		ec.setError(g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up))   // TODO: should be v.KeyK
		ec.setError(g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down)) // TODO: should be v.KeyJ
		ec.setError(g.SetKeybinding(p.name, 'x', gocui.ModNone, v.KeyX))
		ec.setError(g.SetKeybinding(p.name, 'i', gocui.ModNone, v.KeyI))
		ec.setError(g.SetKeybinding(p.name, 'I', gocui.ModNone, v.KeyShiftI))
		ec.setError(g.SetKeybinding(p.name, 'e', gocui.ModNone, v.KeyE))
		ec.setError(g.SetKeybinding(p.name, gocui.KeySpace, gocui.ModNone, v.KeySpace))
	}
	ec.setError(g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, v.KeyEnter))
	ec.setError(g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit)) // TODO: should be v.KeyCtrlC
	return ec.err
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
	t := v.FocusedTask()
	if t == nil {
		log.Printf("focusedTask is nil")
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
		log.Printf("focusedTask is nil")
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
