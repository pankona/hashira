package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	hc "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/daemon"
	"github.com/pankona/hashira/database"
)

func initializeDB() (database.Databaser, error) {
	db := &database.BoltDB{}
	usr, err := user.Current()
	if err != nil {
		return nil, errors.New("failed to current user: " + err.Error())
	}

	configDir := filepath.Join(usr.HomeDir, ".config", "hashira")
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return nil, errors.New("failed to create config directory: " + err.Error())
	}

	err = db.Initialize(filepath.Join(configDir, "db"))
	if err != nil {
		return nil, errors.New("failed to initialize db: " + err.Error())
	}
	return db, nil
}

func main() {
	db, err := initializeDB()
	if err != nil {
		fmt.Printf("failed to initialize DB: %s\n", err.Error())
		os.Exit(1)
	}

	port := 50057
	d := &daemon.Daemon{
		Port: port,
		DB:   db,
	}

	// temporary disable sync feature
	if false {
		const syncPeriod = 10
		go func() {
			for {
				sync(port)
				<-time.After(syncPeriod * time.Second)
			}
		}()
	}

	if err = d.Run(); err != nil {
		fmt.Printf("failed to start hashira daemon: %s\n", err.Error())
		os.Exit(1)
	}
}

const apiServiceURI = "http://localhost:8081/api/v1"

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type UploadRequest struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}

func sync(daemonPort int) {
	cli := &hc.Client{Address: "localhost:" + strconv.Itoa(daemonPort)}
	ts, err := cli.RetrieveAll(context.Background())
	if err != nil {
		log.Println(err)
	}
	p, err := cli.RetrievePriority(context.Background())
	if err != nil {
		log.Println(err)
	}

	tasks := map[string]Task{}
	for k, v := range ts {
		if !v.IsDirty {
			continue
		}
		tasks[k] = Task{
			ID:        v.Id,
			Name:      v.Name,
			Place:     v.Place.String(),
			IsDeleted: v.IsDeleted,
		}
	}
	if len(tasks) == 0 {
		// there's no task to upload
		log.Println("no task to upload to sync")
		return
	}

	log.Printf("%d tasks will upload to sync", len(tasks))

	priority := Priority{}
	for k, v := range p {
		priority[k] = v.Ids
	}

	ur := &UploadRequest{
		Tasks:    tasks,
		Priority: priority,
	}

	body, err := json.Marshal(ur)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiServiceURI+"/upload", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	req.Header["Authorization"] = []string{"09c86189-d824-4617-9675-ed0195e1e233"}

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
