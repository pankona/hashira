package main

import "github.com/jroimartin/gocui"

type Editor struct{}

var HashiraEditor = &Editor{}

func (e *Editor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	gocui.DefaultEditor.Edit(v, key, ch, mod)
}
