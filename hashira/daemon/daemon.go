package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type daemon struct {
	db database.Databaser
}

func (d *daemon) Create(ctx context.Context, cc *service.CommandCreate) (*service.ResultCreate, error) {
	t := &service.Task{
		Name:      cc.GetName(),
		Place:     service.Place_BACKLOG,
		IsDeleted: false,
	}
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, errors.New("failed to create a new task: " + err.Error())
	}

	// specify empty id.
	// expect an id is automatically generated by database itself
	d.db.Save("", buf)
	result := &service.ResultCreate{Task: t}
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

func (d *daemon) Retrieve(ctx context.Context, cr *service.CommandRetrieve) (*service.ResultRetrieve, error) {
	tasks := make([]*service.Task, 0)
	err := d.db.ForEach(func(k, v []byte) error {
		t := &service.Task{}
		err := json.Unmarshal(v, t)
		if err != nil {
			return errors.New("failed to retrieve tasks: " + err.Error())
		}
		tasks = append(tasks, t)
		return nil
	})
	if err != nil {
		return nil, errors.New("failed to retrieve tasks: " + err.Error())
	}
	result := &service.ResultRetrieve{
		Tasks: tasks,
	}
	return result, err
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
	configDir := filepath.Join(usr.HomeDir, ".config", "hashira")
	err = os.Mkdir(configDir, 0755)
	if err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}

	err = db.Initialize(filepath.Join(configDir, "db"))
	if err != nil {
		return errors.New("failed to initialize db: " + err.Error())
	}

	service.RegisterHashiraServer(s, &daemon{db: db})
	return s.Serve(listen)
}
