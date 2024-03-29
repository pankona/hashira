package main

import (
	"context"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pankona/hashira/client"
)

const (
	program = "hashira"
	version = "v0.0.1"
)

func main() {
	c := &client.Client{
		Address: "localhost:50056",
	}
	ctx := context.Background()

	addNewCmd(ctx, c)
	addListCmd(ctx, c)

	kingpin.Version(program + " " + version)
	_ = kingpin.Parse()
}
