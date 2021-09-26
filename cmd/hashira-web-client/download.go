package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	hc "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
)

func download(accesstoken string) {
	log.Println("download started")

	result, err := execDownload(accesstoken)
	if err != nil {
		log.Printf("failed to download task and priority: %v", err)
		return
	}

	log.Printf("%d tasks downloaded", len(result.Tasks))

	cli := &hc.Client{Address: "localhost:" + strconv.Itoa(daemonPort)}
	for _, task := range result.Tasks {
		err = cli.Update(context.Background(), &service.Task{
			Id:        task.ID,
			Name:      task.Name,
			Place:     service.Place(service.Place_value[task.Place]),
			IsDeleted: task.IsDeleted,
			IsDirty:   false,
		})
		if err != nil {
			log.Printf("failed to update task: %v", err)
		}
	}

	priorities := map[string]*service.Priority{}
	for k, v := range result.Priority {
		priorities[k] = &service.Priority{Ids: v}
	}

	_, err = cli.UpdatePriority(context.Background(), priorities)
	if err != nil {
		log.Printf("failed to update priority: %v", err)
	}

	log.Println("download completed")
}

type DownloadResult UploadRequest

func execDownload(accesstoken string) (DownloadResult, error) {
	req, err := http.NewRequest(http.MethodGet, "https://asia-northeast1-hashira-web.cloudfunctions.net/download", nil)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to download tasks and priorities: %w", err)
	}
	defer resp.Body.Close()

	var ret DownloadResult

	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return DownloadResult{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return ret, nil
}
