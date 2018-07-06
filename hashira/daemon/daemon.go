package daemon

import (
	"context"
	"errors"
	"net"
	"os/user"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type daemon struct {
	db database.Databaser
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
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return errors.New("failed to current user: " + err.Error())
	}
	db.Initialize(usr.HomeDir)
	service.RegisterHashiraServer(s, &daemon{db: db})
	return s.Serve(listen)
}
