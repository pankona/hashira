package main

import (
	"context"
	"fmt"
	"log"
	"log/syslog"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/hashira/client"
)

func main() {

	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "hashira-cui")
	log.SetOutput(logger)

	var (
		m  = &Model{}
		ps = &PubSub{}
	)

	// initialize gocui
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// prepare model
	hashirac := &hashirac.Client{
		Address: "localhost:50056",
	}
	m.SetHashiraClient(hashirac)

	// prepare controller
	c := &Ctrl{
		m:   m,
		pub: ps,
	}

	c.Initialize()
	c.SetPublisher(ps)
	// prepare view
	v := &View{}
	err = v.Initialize(g, c)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize view: %s", err.Error()))
	}
	g.SetManager(v)

	err = v.ConfigureKeyBindings(g)
	if err != nil {
		panic(fmt.Sprintf("failed to configure keybindings: %s", err.Error()))
	}

	ps.Subscribe("view", v)

	// retrieve tasks first for initial screen
	err = c.Update(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve initial tasks: %s", err.Error()))
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
