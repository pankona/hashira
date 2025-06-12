package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

// isDaemonRunning checks if the hashira daemon is running on the specified address
func isDaemonRunning(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// startDaemon starts the hashira daemon in background
func startDaemon() error {
	cmd := exec.Command("go", "run", "cmd/hashirad/main.go")
	cmd.Dir = "."
	
	// Start daemon in background
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start daemon: %v", err)
	}
	
	// Detach from the process so it continues running
	go func() {
		_ = cmd.Wait() // Ignore error as this is a background process
	}()
	
	// Wait a bit for daemon to start up
	time.Sleep(5 * time.Second)
	
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

