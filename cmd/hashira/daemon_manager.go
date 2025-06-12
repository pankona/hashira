package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pankona/hashira/daemon"
	"github.com/pankona/hashira/database"
)

// daemonInstance holds the daemon instance for lifecycle management
var daemonInstance *daemon.Daemon

// initializeDB initializes the database for the daemon
func initializeDB() (database.Databaser, error) {
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New("failed to get current user: " + err.Error())
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

// isDaemonRunning checks if the hashira daemon is running on the specified address
func isDaemonRunning(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// startDaemon starts the hashira daemon in background using daemon package directly
func startDaemon() error {
	db, err := initializeDB()
	if err != nil {
		return fmt.Errorf("failed to initialize DB: %v", err)
	}

	daemonInstance = &daemon.Daemon{
		Port: 50057,
		DB:   db,
	}

	// Start daemon in background goroutine
	go func() {
		if err := daemonInstance.Run(); err != nil {
			fmt.Printf("daemon stopped with error: %v\n", err)
		}
	}()

	// Wait a bit for daemon to start up
	time.Sleep(2 * time.Second)
	
	return nil
}

// ensureDaemonRunning checks if daemon is running, and starts it if not
func ensureDaemonRunning(address string) error {
	fmt.Printf("Checking if daemon is running on %s...\n", address)
	if isDaemonRunning(address) {
		fmt.Println("Daemon is already running.")
		return nil
	}
	
	fmt.Println("Hashira daemon is not running. Starting daemon...")
	
	if err := startDaemon(); err != nil {
		return fmt.Errorf("failed to start daemon: %v", err)
	}
	
	fmt.Println("Waiting for daemon to start...")
	// Verify daemon started successfully
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if isDaemonRunning(address) {
			fmt.Println("Daemon started successfully.")
			return nil
		}
		fmt.Printf("Waiting... (%d/%d)\n", i+1, maxRetries)
		time.Sleep(1 * time.Second)
	}
	
	return fmt.Errorf("daemon failed to start within %d seconds", maxRetries)
}

// stopDaemon stops the daemon if it's running
func stopDaemon() {
	if daemonInstance != nil {
		daemonInstance.Stop()
		daemonInstance = nil
	}
}

