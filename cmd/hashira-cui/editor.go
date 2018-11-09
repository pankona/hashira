package main

import (
	"github.com/jroimartin/gocui"
)

type Editor struct{}

var HashiraEditor = &Editor{}

func (e *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch key {
	case gocui.KeyCtrlH:
		v.MoveCursor(+1, 0, false)
	case gocui.KeyCtrlL:
		v.MoveCursor(-1, 0, false)
	}
	gocui.DefaultEditor.Edit(v, key, ch, mod)
}
