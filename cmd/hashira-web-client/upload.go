package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	hc "github.com/pankona/hashira/client"
)

func upload(accesstoken string) {
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

	if len(tasks) == 0 {
		// there's no task to upload
		log.Println("no task to upload to sync")
		return
	}

	priority := Priority{}
	for k, v := range p {
		priority[k] = v.Ids
	}

	err = execUpload(accesstoken, tasks, priority)
	if err != nil {
		log.Printf("failed to upload sync: %v", err)
		return
	}
}

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type UploadRequest struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}

func execUpload(accesstoken string, tasks map[string]Task, priority Priority) error {
	log.Printf("%d tasks will upload to sync", len(tasks))

	ur := &UploadRequest{
		Tasks:    tasks,
		Priority: priority,
	}

	body, err := json.Marshal(ur)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://asia-northeast1-hashira-web.cloudfunctions.net/upload", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}