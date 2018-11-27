package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
	"gopkg.in/alecthomas/kingpin.v2"
)

func addNewCmd(ctx context.Context, c *client.Client) {
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
}

func addListCmd(ctx context.Context, c *client.Client) {
	listCmd := kingpin.Command(
		"list",
		"show list of tasks")
	_ = listCmd.Action(func(pc *kingpin.ParseContext) error {
		return list(ctx, c)
	})
}

func create(ctx context.Context, c *client.Client, name string) error {
	t := &service.Task{
		Name:  name,
		Place: service.Place_BACKLOG,
	}
	err := c.Create(ctx, t)
	if err != nil {
		return errors.New("failed to create a new task: " + err.Error())
	}
	return nil
}

func list(ctx context.Context, c *client.Client) error {
	tasks, err := c.Retrieve(ctx)
	if err != nil {
		return errors.New("failed to create a new task: " + err.Error())
	}
	for _, v := range tasks {
		id, err := strconv.Atoi(v.Id)
		if err != nil {
			continue
		}
		fmt.Printf("[%04d]\t%s\t%v\t%v\n", id, v.Name, v.Place, v.IsDeleted)
	}
	return nil
}
