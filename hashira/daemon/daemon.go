package daemon

import (
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

type Daemon struct {
	DB database.Databaser
}

func NewDemon() *Daemon {
	return &Daemon{}
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
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return errors.New("failed to current user: " + err.Error())
	}
	configDir := filepath.Join(usr.HomeDir, ".config", "hashira")
	err = os.Mkdir(configDir, 0700)
	if err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}

	err = db.Initialize(filepath.Join(configDir, "db"))
	if err != nil {
		return errors.New("failed to initialize db: " + err.Error())
	}

	service.RegisterHashiraServer(s, &Daemon{DB: db})
	return s.Serve(listen)
}
