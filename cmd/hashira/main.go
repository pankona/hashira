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
	address := "localhost:50057"
	c := &client.Client{
		Address: address,
	}
	ctx := context.Background()

	addNewCmd(ctx, c, address)
	addListCmd(ctx, c, address)

	kingpin.Version(program + " " + version)
	_ = kingpin.Parse()
}
