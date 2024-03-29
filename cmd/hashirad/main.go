package main

import (
	"context"
	"errors"
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
	db, err := initializeDB()
	if err != nil {
		fmt.Printf("failed to initialize DB: %s\n", err.Error())
		os.Exit(1)
	}

	port := 50057
	d := &daemon.Daemon{
		Port: port,
		DB:   db,
	}

	accesstoken, ok := os.LookupEnv("HASHIRA_ACCESS_TOKEN")
	if ok {
		if err := startSync(context.Background(), port, accesstoken); err != nil {
			log.Printf("failed to start synchronization: %v", err)
		}
	}

	if err = d.Run(); err != nil {
		fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
		os.Exit(1)
	}
}

func startSync(ctx context.Context, daemonPort int, accesstoken string) error {
	fmt.Printf("start synchronization...\n")

	sc := syncutil.Client{DaemonPort: daemonPort}
	err := sc.TestAccessToken(accesstoken)
	if err != nil {
		return fmt.Errorf("HASHIRA_ACCESSTOKEN is invalid. Synchronization is not started: %w", err)
	}
	fmt.Printf("HASHIRA_ACCESSTOKEN is valid. hashira-web will work!\n")

	const syncInterval = 10 * time.Minute

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				if err := sc.Upload(accesstoken, syncutil.UploadDirtyOnly); err != nil {
					log.Printf("failed to upload: %v", err)
				}
				if err := sc.Download(accesstoken); err != nil {
					log.Printf("failed to download: %v", err)
				}
				<-time.After(syncInterval)
			}
		}
	}()

	return nil
}
