package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
)

func withClient(f func(service.HashiraClient) error) error {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
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
func Create(ctx context.Context, taskName string) error {
	return withClient(
		func(hc service.HashiraClient) error {
			cc := &service.CommandCreate{
				Name:  taskName,
				Place: service.Place_BACKLOG,
			}
			result, err := hc.Create(ctx, cc)
			if err != nil {
				return errors.New("Create failed: " + err.Error())
			}
			result.ProtoMessage()
			return nil
		})
}

// Retrieve retrieves all tasks
func Retrieve(ctx context.Context) ([]*service.Task, error) {
	var tasks []*service.Task
	err := withClient(
		func(hc service.HashiraClient) error {
			cr := &service.CommandRetrieve{}
			result, err := hc.Retrieve(ctx, cr)
			if err != nil {
				return errors.New("Retrieve failed: " + err.Error())
			}
			tasks = result.Tasks
			return nil
		})
	return tasks, err
}
