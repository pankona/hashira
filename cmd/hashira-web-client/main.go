package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pankona/hashira/daemon"
	"github.com/pankona/hashira/database"
)

func main() {
	var (
		flagSync     bool
		flagUpload   bool
		flagDownload bool
		flagTest     bool
	)

	flag.BoolVar(&flagSync, "sync", false, "sync (download and upload) tasks and priorities with hashira-web")
	flag.BoolVar(&flagUpload, "upload", false, "upload tasks and priorities to hashira-web")
	flag.BoolVar(&flagDownload, "download", false, "download tasks and priorities from hashira-web")
	flag.BoolVar(&flagTest, "test", false, "test the hashira-web works")
	flag.Parse()

	accesstoken := os.Getenv("HASHIRA_ACCESS_TOKEN")
	if accesstoken == "" {
		log.Printf("Please specify environment variable HASHIRA_ACCESS_TOKEN. Abort.")
		os.Exit(1)
	}

	done := make(chan struct{})
	go func() {
		launchHashirad(done)
	}()

	time.Sleep(1 * time.Second)

	switch {
	case flagSync:
		download(accesstoken)
		upload(accesstoken, UploadAll)
	case flagUpload:
		upload(accesstoken, UploadDirtyOnly)
	case flagDownload:
		download(accesstoken)
	case flagTest:
		fallthrough
	default:
		testAccessToken(accesstoken)
	}

	done <- struct{}{}
}

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

const daemonPort = 50058

func launchHashirad(done <-chan struct{}) {
	db, err := initializeDB()
	if err != nil {
		fmt.Printf("failed to initialize DB: %s\n", err.Error())
		os.Exit(1)
	}

	d := &daemon.Daemon{
		Port: daemonPort,
		DB:   db,
	}

	go func() {
		<-done
		d.Stop()
	}()

	if err = d.Run(); err != nil {
		fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
		os.Exit(1)
	}
}
