package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

func main() {

	var (
		m  = &Model{}
		ps = &PubSub{}
	)

	// prepare model
	hashirac := &hashirac.Client{
		Address: "localhost:50056",
	}
	m.SetHashiraClient(hashirac)
	m.SetPublisher(ps)

	// prepare view
	v := &View{}
	err := v.Initialize()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize view: %s", err.Error()))
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManager(v)

	ps.Subscribe("view", v)

	// prepare controller
	c := &Ctrl{
		m: m,
		g: g,
	}

	err = c.ConfigureKeyBindings(g)
	if err != nil {
		panic(fmt.Sprintf("failed to configure keybindings: %s", err.Error()))
	}

	// retrieve tasks first for initial screen
	err = c.Update(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve initial tasks: %s", err.Error()))
	}

	// start to run main loop
	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
