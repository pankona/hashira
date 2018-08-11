package main

import (
	"context"
	"errors"

	"fmt"
	"strconv"

	"github.com/pankona/hashira/hashira/client"
)

func create(ctx context.Context, c *client.Client, name string) error {
	err := c.Create(ctx, name)
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
