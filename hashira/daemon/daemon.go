package daemon

import (
	"errors"
	"net"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Daemon is a structure that implement hashira service
type Daemon struct {
	DB database.Databaser
}

// Run starts hashira daemon (as gRPC server)
func (d *Daemon) Run() error {
	port := ":50056" // TODO: specify port number via function argument
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return errors.New("gRPC server failed to listen port " + port + ": " + err.Error())
	}
	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterHashiraServer(s, d)
	return s.Serve(listen)
}
