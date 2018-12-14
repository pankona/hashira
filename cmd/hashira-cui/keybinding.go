package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/jesseduffield/gocui"
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

		ec.setError(g.SetKeybinding(p.name, 'h', gocui.ModNone, v.KeyLeft))
		ec.setError(g.SetKeybinding(p.name, gocui.KeyArrowLeft, gocui.ModNone, v.KeyLeft))
		ec.setError(g.SetKeybinding(p.name, 'l', gocui.ModNone, v.KeyRight))
		ec.setError(g.SetKeybinding(p.name, gocui.KeyArrowRight, gocui.ModNone, v.KeyRight))

		ec.setError(g.SetKeybinding(p.name, 'k', gocui.ModNone, v.Up))
		ec.setError(g.SetKeybinding(p.name, gocui.KeyArrowUp, gocui.ModNone, v.Up))
		ec.setError(g.SetKeybinding(p.name, 'j', gocui.ModNone, v.Down))
		ec.setError(g.SetKeybinding(p.name, gocui.KeyArrowDown, gocui.ModNone, v.Down))

		ec.setError(g.SetKeybinding(p.name, 'x', gocui.ModNone, v.KeyX))
		ec.setError(g.SetKeybinding(p.name, 'i', gocui.ModNone, v.KeyI))
		ec.setError(g.SetKeybinding(p.name, 'I', gocui.ModNone, v.KeyShiftI))
		ec.setError(g.SetKeybinding(p.name, 'e', gocui.ModNone, v.KeyE))
		ec.setError(g.SetKeybinding(p.name, gocui.KeySpace, gocui.ModNone, v.KeySpace))
	}
	ec.setError(g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, v.KeyEnter))
	ec.setError(g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, v.KeyEsc))
	ec.setError(g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, v.Quit))
	return ec.err
}

// KeyLeft reacts for "h" and left arrow
func (v *View) KeyLeft(g *gocui.Gui, _ *gocui.View) error {
	return v.Left()
}

// KeyRight reacts for "l" and right arrow
func (v *View) KeyRight(g *gocui.Gui, _ *gocui.View) error {
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

// KeyEsc reacts for "Esc"
func (v *View) KeyEsc(g *gocui.Gui, gv *gocui.View) error {
	if gv.Name() == "input" {
		return v.hideInput(g, gv)
	}
	return nil
}
