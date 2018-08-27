package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	ps := &PubSub{}

	view := &view{}
	ps.Subscribe("view", view)

	g.SetManager(view)

	ctrl := NewCtrl()

	err = ctrl.ConfigureKeyBindings(g)
	if err != nil {
		panic(fmt.Sprintf("failed to configure keybindings: %s", err.Error()))
	}

	hashirac := &hashirac.Client{
		Address: "localhost:50056",
	}
	ctrl.SetHashiraClient(hashirac)

	ctrl.SetPublisher(ps)

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
