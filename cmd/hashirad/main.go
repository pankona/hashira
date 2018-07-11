package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/hashira/daemon"
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

	d := &daemon.Daemon{
		Port: 50056,
		DB:   db,
	}
	if err = d.Run(); err != nil {
		fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
		os.Exit(1)
	}
}
