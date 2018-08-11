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

	newCmd := kingpin.Command(
		"new",
		"add new task with specified task name")
	name := newCmd.Arg(
		"name",
		"name of task which is newly added").
		Required().String()
	_ = newCmd.Action(func(pc *kingpin.ParseContext) error {
		return create(ctx, c, *name)
	})

	listCmd := kingpin.Command(
		"list",
		"show list of tasks")
	_ = listCmd.Action(func(pc *kingpin.ParseContext) error {
		return list(ctx, c)
	})

	kingpin.Parse()
}
