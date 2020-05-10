package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/pankona/hashira/service"
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

	accesstoken := os.Getenv("HASHIRA_ACCESSTOKEN")

	if len(accesstoken) != 0 {
		const syncPeriod = 10

		go func() {
			<-time.After(syncPeriod * time.Second)
			initialSync(port, accesstoken)

			// TODO: don't enter infinite loop if accesstoken is invalid
			for {
				<-time.After(syncPeriod * time.Second)
				sync(port, accesstoken)
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

func initialSync(daemonPort int, accesstoken string) {
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
		tasks[k] = Task{
			ID:        v.Id,
			Name:      v.Name,
			Place:     v.Place.String(),
			IsDeleted: v.IsDeleted,
		}
	}

	priority := Priority{}
	if len(tasks) != 0 {
		for k, v := range p {
			priority[k] = v.Ids
		}

		err = upload(accesstoken, tasks, priority)
		if err != nil {
			log.Printf("failed to upload sync: %v", err)
			return
		}
	}

	tasks, priority, err = download(accesstoken)
	if err != nil {
		log.Printf("failed to download: %v", err)
		return
	}

	for _, v := range tasks {
		err = cli.Update(context.Background(), &service.Task{
			Id:        v.ID,
			Name:      v.Name,
			Place:     service.Place(service.Place_value[v.Place]),
			IsDeleted: v.IsDeleted,
			IsDirty:   false,
		})
		if err != nil {
			log.Println(err)
		}
	}

	p = map[string]*service.Priority{}
	for k, v := range priority {
		p[k] = &service.Priority{Ids: v}
	}

	log.Printf("downloaded priorities:")
	for k, v := range p {
		fmt.Printf("%s:\n", k)
		for _, id := range v.GetIds() {
			fmt.Printf("%s, ", tasks[id].Name)
		}
		fmt.Printf("\n")
	}
	_, err = cli.UpdatePriority(context.Background(), p)
	if err != nil {
		log.Println(err)
	}
}

func sync(daemonPort int, accesstoken string) {
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

	priority := Priority{}
	for k, v := range p {
		priority[k] = v.Ids
	}

	err = upload(accesstoken, tasks, priority)
	if err != nil {
		log.Printf("failed to upload sync: %v", err)
		return
	}

	tasks, priority, err = download(accesstoken)
	if err != nil {
		log.Printf("failed to download: %v", err)
		return
	}

	for _, v := range tasks {
		err = cli.Update(context.Background(), &service.Task{
			Id:        v.ID,
			Name:      v.Name,
			Place:     service.Place(service.Place_value[v.Place]),
			IsDeleted: v.IsDeleted,
			IsDirty:   false,
		})
		if err != nil {
			log.Println(err)
		}
	}

	p = map[string]*service.Priority{}
	for k, v := range priority {
		p[k] = &service.Priority{Ids: v}
	}

	log.Printf("downloaded priorities:")
	for k, v := range p {
		fmt.Printf("%s:\n", k)
		for _, id := range v.GetIds() {
			fmt.Printf("%s, ", tasks[id].Name)
		}
		fmt.Printf("\n")
	}
	_, err = cli.UpdatePriority(context.Background(), p)
	if err != nil {
		log.Println(err)
	}
}

func upload(accesstoken string, tasks map[string]Task, priority Priority) error {
	log.Printf("%d tasks will upload to sync", len(tasks))

	ur := &UploadRequest{
		Tasks:    tasks,
		Priority: priority,
	}

	log.Printf("uploading priorities:")
	for k, v := range priority {
		fmt.Printf("%s:\n", k)
		for _, id := range v {
			fmt.Printf("%s, ", tasks[id].Name)
		}
		fmt.Printf("\n")
	}

	body, err := json.Marshal(ur)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiServiceURI+"/upload", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	req.Header["Authorization"] = []string{accesstoken}

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type DownloadResponse UploadRequest

func download(accesstoken string) (map[string]Task, Priority, error) {
	req, err := http.NewRequest(http.MethodGet, apiServiceURI+"/download", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header["Authorization"] = []string{accesstoken}

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	dr := DownloadResponse{}
	err = json.Unmarshal(buf, &dr)
	if err != nil {
		return nil, nil, err
	}

	return dr.Tasks, dr.Priority, nil
}
