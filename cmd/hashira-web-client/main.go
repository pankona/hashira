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
	"github.com/pankona/hashira/sync/syncutil"
)

var (
	Version  = "unset"
	Revision = "unset"
)

const daemonPort = 50058

func main() {
	var (
		flagSync     bool
		flagUpload   bool
		flagDownload bool
		flagTest     bool
		flagVersion  bool
	)

	flag.BoolVar(&flagVersion, "version", false, "show version")
	flag.BoolVar(&flagSync, "sync", false, "sync (download and upload) tasks and priorities with hashira-web")
	flag.BoolVar(&flagUpload, "upload", false, "upload tasks and priorities to hashira-web")
	flag.BoolVar(&flagDownload, "download", false, "download tasks and priorities from hashira-web")
	flag.BoolVar(&flagTest, "test", false, "test the hashira-web works")
	flag.Parse()

	if flagVersion {
		fmt.Printf("hashira-web-client version: %s, Revision: %s\n", Version, Revision)
		return
	}

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

	sc := syncutil.Client{DaemonPort: daemonPort}

	switch {
	case flagSync:
		if err := sc.Download(accesstoken); err != nil {
			log.Printf("failed to download: %v", err)
		}
		log.Printf("download completed")
		if err := sc.Upload(accesstoken, syncutil.UploadAll); err != nil {
			log.Printf("failed to upload: %v", err)
		}
		log.Printf("upload completed")
	case flagUpload:
		if err := sc.Upload(accesstoken, syncutil.UploadDirtyOnly); err != nil {
			log.Printf("failed to upload: %v", err)
		}
		log.Printf("upload completed")
	case flagDownload:
		if err := sc.Download(accesstoken); err != nil {
			log.Printf("failed to download: %v", err)
		}
		log.Printf("download completed")
	case flagTest:
		fallthrough
	default:
		if err := sc.TestAccessToken(accesstoken); err != nil {
			log.Printf("test access token failed: %v", err)
		} else {
			log.Println("The accesstoken is valid. hashira-web will work!")
		}
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
