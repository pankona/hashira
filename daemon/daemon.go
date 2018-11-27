package daemon

import (
	"errors"
	"net"
	"strconv"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Daemon is a structure that implement hashira service
type Daemon struct {
	Port   int
	DB     database.Databaser
	server *grpc.Server
}

// Run starts hashira daemon (as gRPC server)
func (d *Daemon) Run() error {
	p := ":" + strconv.Itoa(d.Port)
	listen, err := net.Listen("tcp", p)
	if err != nil {
		return errors.New("gRPC server failed to listen [" + p + "]: " + err.Error())
	}
	s := grpc.NewServer()
	reflection.Register(s)
	service.RegisterHashiraServer(s, d)

	d.server = s

	return d.server.Serve(listen)
}

// Stop stops hashira daemon
func (d *Daemon) Stop() {
	d.server.Stop()
}
