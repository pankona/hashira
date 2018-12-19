package main

import (
	runewidth "github.com/mattn/go-runewidth"
	"github.com/pankona/gocui"
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
		_, oy := v.Origin()
		v.SetOrigin(0, oy)

		maxX, _ := v.Size()
		v.MoveCursor(-maxX, 0, false)
	case gocui.KeyCtrlE:
		// move to end of line
		_, oy := v.Origin()
		v.SetOrigin(0, oy)

		bufLen := runewidth.StringWidth(v.Buffer())
		cx, _ := v.Cursor()
		v.MoveCursor(+bufLen-cx, 0, false)
	case gocui.KeyCtrlD:
		v.EditDelete(false)
	default:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
