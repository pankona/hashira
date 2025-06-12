package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
)

func addNewCmd(ctx context.Context, c *client.Client, address string) {
	newCmd := kingpin.Command(
		"new",
		"add new task with specified task name")
	name := newCmd.Arg(
		"name",
		"name of task which is newly added").
		Required().String()
	_ = newCmd.Action(func(pc *kingpin.ParseContext) error {
		if err := ensureDaemonRunning(address); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return create(ctx, c, *name)
	})
}

func addListCmd(ctx context.Context, c *client.Client, address string) {
	listCmd := kingpin.Command(
		"list",
		"show list of tasks")
	filter := listCmd.Flag(
		"filter",
		"filter tasks by name (case-insensitive substring match)").
		String()
	_ = listCmd.Action(func(pc *kingpin.ParseContext) error {
		if err := ensureDaemonRunning(address); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
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
	
	// Normalize filter once before the loop for better performance
	normalizedFilter := strings.ToLower(filter)
	
	filteredTasks := make(map[string]*service.Task)
	for k, v := range tasks {
		if normalizedFilter == "" || strings.Contains(strings.ToLower(v.Name), normalizedFilter) {
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
