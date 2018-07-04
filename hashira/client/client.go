package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
)

type client struct {
	hc service.HashiraClient
}

func withClient(f func(service.HashiraClient) error) error {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		// TODO: error handling
	}
	defer func() {
		_ = conn.Close()
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
