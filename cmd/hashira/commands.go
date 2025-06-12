package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
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
	statuses := listCmd.Flag(
		"status",
		"comma-separated list of statuses to show (BACKLOG,TODO,DOING,DONE). Default: BACKLOG,TODO,DOING").
		String()
	_ = listCmd.Action(func(pc *kingpin.ParseContext) error {
		if err := ensureDaemonRunning(address); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return list(ctx, c, *filter, *statuses)
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

func list(ctx context.Context, c *client.Client, filter string, statusFilter string) error {
	tasks, err := c.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %v", err)
	}

	// Parse status filter - default excludes DONE
	allowedStatuses := map[service.Place]bool{
		service.Place_BACKLOG: true,
		service.Place_TODO:    true,
		service.Place_DOING:   true,
		service.Place_DONE:    false,
	}

	if statusFilter != "" {
		// Reset all to false if custom filter provided
		for k := range allowedStatuses {
			allowedStatuses[k] = false
		}

		statusList := strings.Split(strings.ToUpper(statusFilter), ",")
		for _, status := range statusList {
			status = strings.TrimSpace(status)
			switch status {
			case "BACKLOG":
				allowedStatuses[service.Place_BACKLOG] = true
			case "TODO":
				allowedStatuses[service.Place_TODO] = true
			case "DOING":
				allowedStatuses[service.Place_DOING] = true
			case "DONE":
				allowedStatuses[service.Place_DONE] = true
			}
		}
	}

	// Normalize filter once before the loop for better performance
	normalizedFilter := strings.ToLower(filter)

	var filteredTasks []*service.Task
	for _, v := range tasks {
		// Filter by status
		if !allowedStatuses[v.Place] {
			continue
		}

		// Filter by name
		if normalizedFilter != "" && !strings.Contains(strings.ToLower(v.Name), normalizedFilter) {
			continue
		}

		filteredTasks = append(filteredTasks, v)
	}

	if len(filteredTasks) == 0 {
		if filter != "" {
			fmt.Printf("No tasks found matching filter: %s\n", filter)
		} else {
			fmt.Println("No tasks found.")
		}
		return nil
	}

	// Sort by status order: BACKLOG -> TODO -> DOING -> DONE
	sort.Slice(filteredTasks, func(i, j int) bool {
		return filteredTasks[i].Place < filteredTasks[j].Place
	})

	for _, v := range filteredTasks {
		fmt.Printf("[%s]\t%s\t%v\t%v\n", v.Id, v.Name, v.Place, v.IsDeleted)
	}
	return nil
}
