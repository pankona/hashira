package daemon

import (
	"context"
	"errors"
	"net"

	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type daemon struct {
}

func (d *daemon) Create(ctx context.Context, cc *service.CommandCreate) (*service.ResultCreate, error) {
	// TODO: implement
	t := &service.Task{
		Id:        "TODO: generate proper id",
		Name:      cc.GetName(),
		Place:     service.Place_BACKLOG,
		IsDeleted: false,
	}
	result := &service.ResultCreate{t}
	return result, nil
}

func (d *daemon) Update(context.Context, *service.CommandUpdate) (*service.ResultUpdate, error) {
	// TODO: implement
	return nil, nil
}

func (d *daemon) Delete(context.Context, *service.CommandDelete) (*service.ResultDelete, error) {
	// TODO: implement
	return nil, nil
}

func (d *daemon) Retrieve(context.Context, *service.CommandRetrieve) (*service.ResultRetrieve, error) {
	// TODO: implement
	return nil, nil
}

func Run() error {
	port := ":50056" // TODO: specify port number via function argument
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return errors.New("gRPC server failed to listen port " + port + ": " + err.Error())
	}
	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterHashiraServer(s, &daemon{})
	return s.Serve(listen)
}
