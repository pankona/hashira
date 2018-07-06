package main

import (
	"fmt"
	"os"

	"github.com/pankona/hashira/hashira/daemon"
)

func main() {
	err := daemon.Run()
	if err != nil {
		fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
		os.Exit(1)
	}
}
