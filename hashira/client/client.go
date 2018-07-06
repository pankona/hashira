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

func Create(ctx context.Context) error {
	return withClient(
		func(hc service.HashiraClient) error {
			cc := &service.CommandCreate{
				Name:  "test",
				Place: service.Place_BACKLOG,
			}
			result, err := hc.Create(ctx, cc)
			if err != nil {
				return errors.New("Create failed: " + err.Error())
			}
			result.ProtoMessage()
			fmt.Println(result.String())
			return nil
		})
}
