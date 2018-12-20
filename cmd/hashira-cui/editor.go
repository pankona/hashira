package main

import (
	"log"

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
		err := v.SetOrigin(0, oy)
		if err != nil {
			log.Printf("failed to set origin: %v", err)
			return
		}

		maxX, _ := v.Size()
		v.MoveCursor(-maxX, 0, false)
	case gocui.KeyCtrlE:
		// move to end of line
		_, oy := v.Origin()
		err := v.SetOrigin(0, oy)
		if err != nil {
			log.Printf("failed to set origin: %v", err)
			return
		}

		bufLen := runewidth.StringWidth(v.Buffer())
		cx, _ := v.Cursor()
		v.MoveCursor(+bufLen-cx, 0, false)
	case gocui.KeyCtrlD:
		v.EditDelete(false)
	default:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
