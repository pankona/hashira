package main

import (
	"github.com/jesseduffield/gocui"
)

type editor struct{}

var hashiraEditor = &editor{}

func (e *editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch key {
	case gocui.KeyCtrlF:
		v.MoveCursor(+1, 0, false)
	case gocui.KeyCtrlB:
		v.MoveCursor(-1, 0, false)
	case gocui.KeyCtrlA:
		// move to start of line
		maxX, _ := v.Size()
		v.MoveCursor(-maxX, 0, false)
	case gocui.KeyCtrlE:
		// move to end of line
		bufLen := len(v.Buffer())
		cx, _ := v.Cursor()
		v.MoveCursor(+bufLen-cx-1, 0, false)
	case gocui.KeyCtrlD:
		v.EditDelete(false)
	default:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
