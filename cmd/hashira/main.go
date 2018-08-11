package main

import (
	"context"

	"github.com/pankona/hashira/hashira/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	c := &client.Client{
		Address: "localhost:50056",
	}
	ctx := context.Background()

	addNewCmd(ctx, c)
	addListCmd(ctx, c)

	kingpin.Parse()
}
