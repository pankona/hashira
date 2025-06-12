package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
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
	filter := listCmd.Flag(
		"filter",
		"filter tasks by name (case-insensitive substring match)").
		String()
	_ = listCmd.Action(func(pc *kingpin.ParseContext) error {
		return list(ctx, c, *filter)
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

func list(ctx context.Context, c *client.Client, filter string) error {
	tasks, err := c.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %v", err)
	}
	
	filteredTasks := make(map[string]*service.Task)
	for k, v := range tasks {
		if filter == "" || strings.Contains(strings.ToLower(v.Name), strings.ToLower(filter)) {
			filteredTasks[k] = v
		}
	}
	
	if len(filteredTasks) == 0 {
		if filter != "" {
			fmt.Printf("No tasks found matching filter: %s\n", filter)
		} else {
			fmt.Println("No tasks found.")
		}
		return nil
	}
	
	for _, v := range filteredTasks {
		fmt.Printf("[%s]\t%s\t%v\t%v\n", v.Id, v.Name, v.Place, v.IsDeleted)
	}
	return nil
}
