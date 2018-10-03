package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
)

// Client is a hashira client structure
type Client struct {
	Address string
}

func (c *Client) withClient(f func(service.HashiraClient) error) error {
	conn, err := grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return errors.New("failed to Dial: " + err.Error())
	}
	defer func() {
		e := conn.Close()
		if e != nil {
			fmt.Printf("failed to close connection: %s\n", e.Error())
		}
	}()

	return f(service.NewHashiraClient(conn))
}

// Create creates new task
func (c *Client) Create(ctx context.Context, task *service.Task) error {
	return c.withClient(
		func(hc service.HashiraClient) error {
			com := &service.CommandCreate{
				Task: task,
			}
			result, err := hc.Create(ctx, com)
			if err != nil {
				return errors.New("Create failed: " + err.Error())
			}
			result.ProtoMessage()
			return nil
		})
}

// Update updates an existing task
func (c *Client) Update(ctx context.Context, task *service.Task) error {
	return c.withClient(
		func(hc service.HashiraClient) error {
			com := &service.CommandUpdate{
				Task: task,
			}

			_, err := hc.Update(context.Background(), com)
			if err != nil {
				return fmt.Errorf("update a task failed: %s", err.Error())
			}

			return nil
		})
}

// Delete marks specified task as deleted
func (c *Client) Delete(ctx context.Context, id string) error {
	return c.withClient(
		func(hc service.HashiraClient) error {
			com := &service.CommandDelete{
				Id: id,
			}
			result, err := hc.Delete(ctx, com)
			if err != nil {
				return errors.New("Delete failed: " + err.Error())
			}
			result.ProtoMessage()
			return nil
		})
}

// Retrieve retrieves all tasks
func (c *Client) Retrieve(ctx context.Context) ([]*service.Task, error) {
	var tasks []*service.Task
	err := c.withClient(
		func(hc service.HashiraClient) error {
			com := &service.CommandRetrieve{}
			result, err := hc.Retrieve(ctx, com)
			if err != nil {
				return errors.New("Retrieve failed: " + err.Error())
			}
			tasks = result.Tasks
			return nil
		})
	return tasks, err
}

// UpdatePriority updates tasks' priorities
func (c *Client) UpdatePriority(ctx context.Context, priorities []*service.Priority) ([]*service.Priority, error) {
	var ret []*service.Priority

	err := c.withClient(func(hc service.HashiraClient) error {
		com := &service.CommandUpdatePriority{
			Priorities: priorities,
		}
		p, err := hc.UpdatePriority(ctx, com)
		if err != nil {
			return errors.New("UpdatePriority failed: " + err.Error())
		}
		ret = p.Priorities
		return nil
	})

	return ret, err
}

// RetrievePriority retrieves tasks' priorities
func (c *Client) RetrievePriority(ctx context.Context) ([]*service.Priority, error) {
	var ret []*service.Priority

	err := c.withClient(func(hc service.HashiraClient) error {
		com := &service.CommandRetrievePriority{}
		p, err := hc.RetrievePriority(ctx, com)
		if err != nil {
			return errors.New("RetrievePriority failed: " + err.Error())
		}
		ret = p.Priorities
		return nil
	})

	log.Printf("client retrieves priority: %v", ret)

	return ret, err
}
