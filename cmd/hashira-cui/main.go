package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jroimartin/gocui"
	hashirac "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/daemon"
	"github.com/pankona/hashira/database"
)

func initializeDB() (database.Databaser, error) {
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New("failed to current user: " + err.Error())
	}

	configDir := filepath.Join(usr.HomeDir, ".config", "hashira")
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return nil, errors.New("failed to create config directory: " + err.Error())
	}

	err = db.Initialize(filepath.Join(configDir, "db"))
	if err != nil {
		return nil, errors.New("failed to initialize db: " + err.Error())
	}
	return db, nil
}

func main() {
	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "hashira-cui")
	log.SetOutput(logger)

	db, err := initializeDB()
	if err != nil {
		os.Exit(1)
	}

	d := &daemon.Daemon{
		Port: 50056,
		DB:   db,
	}

	go func() {
		if err = d.Run(); err != nil {
			fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
			os.Exit(1)
		}
	}()
	defer func() {
		d.Stop()
	}()

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
