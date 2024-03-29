package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pankona/gocui"
	hashirac "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/daemon"
	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/sync/syncutil"
	"github.com/pankona/hashira/xdg"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	var (
		flagVersion bool
	)
	flag.BoolVar(&flagVersion, "version", false, "show version")
	flag.Parse()

	if flagVersion {
		fmt.Printf("hashira-cui version: %s, Revision: %s\n", Version, Revision)
		return
	}

	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "hashira-cui")
	if err != nil {
		log.Printf("failed to connect to logger but continue to work: %v", err)
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(logger)
	}

	db, err := initializeDB()
	if err != nil {
		os.Exit(1)
	}

	const daemonPort = 50056

	d := &daemon.Daemon{
		Port: daemonPort,
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

	// initialize gocui
	// specify false means: supportOverlaps = false
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// prepare model
	hashirac := &hashirac.Client{Address: fmt.Sprintf("localhost:%d", daemonPort)}
	syncclient := &syncutil.Client{DaemonPort: daemonPort}
	m := NewModel(hashirac, syncclient)

	// Start synchronization with cloud if HASHIRA_ACCESS_TOKEN is set
	accesstoken, ok := os.LookupEnv("HASHIRA_ACCESS_TOKEN")

	if ok {
		go func() {
			sc := syncutil.Client{DaemonPort: daemonPort}
			err := sc.TestAccessToken(accesstoken)
			if err != nil {
				log.Printf("HASHIRA_ACCESSTOKEN is invalid. Synchronization is not started: %v", err)
			}
			m.SetAccessToken(accesstoken)

			go func() {
				// sync on launch immediately
				if err := m.SyncNow(context.Background()); err != nil {
					log.Printf("failed to sync: %v", err)
				}

				// start polling
				for {
					m.NotifySync()
					<-time.After(2 * time.Minute)
				}
			}()

			if err := m.SyncOnNotify(context.Background()); err != nil {
				log.Printf("sync on notify finished: %v", err)
			}
		}()
	}

	// prepare controller
	ps := &PubSub{}
	m.SetPublisher(ps)

	c := &Ctrl{
		m:   m,
		pub: ps,
	}

	c.Initialize()
	c.SetPublisher(ps)
	// prepare view
	v := &View{}
	v.Initialize(g, c)
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

func initializeDB() (database.Databaser, error) {
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New("failed to current user: " + err.Error())
	}

	legacyDataDir := filepath.Join(usr.HomeDir, ".config", "hashira")
	isLegacyDataDirExist := func() bool {
		legacyInfo, err := os.Stat(legacyDataDir)
		if err != nil && os.IsNotExist(err) {
			return false
		}
		return legacyInfo.IsDir()
	}
	dataDir := legacyDataDir
	if !isLegacyDataDirExist() {
		x := &xdg.Xdg{User: *usr}
		dataDir = filepath.Join(x.DataHome(), "hashira")
		err = os.MkdirAll(dataDir, 0700)
		if err != nil {
			return nil, errors.New("failed to create data directory: " + err.Error())
		}
	}

	err = db.Initialize(filepath.Join(dataDir, "db"))
	if err != nil {
		return nil, errors.New("failed to initialize db: " + err.Error())
	}
	return db, nil
}
