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

	newPriorities := map[string]*service.Priority{}
	for k, v := range result.Priority {
		newPriorities[k] = &service.Priority{Ids: v}
	}

	oldPriorities, err := cli.RetrievePriority(context.Background())
	if err != nil {
		log.Printf("failed to retrieve old priorities: %v", err)
	}

	priorities := mergePriorities(newPriorities, oldPriorities)

	_, err = cli.UpdatePriority(context.Background(), priorities)
	if err != nil {
		log.Printf("failed to update priority: %v", err)
	}

	log.Println("download completed")
}

func mergePriorities(newPriorities, oldPriorities map[string]*service.Priority) map[string]*service.Priority {
	ret := map[string]*service.Priority{}
	for k, oldPriority := range oldPriorities {
		ret[k] = &service.Priority{
			Ids: append(newPriorities[k].Ids, oldPriority.Ids...),
		}
		ret[k].Ids = unique(ret[k].Ids)
	}
	return ret
}

func unique(ss []string) []string {
	keys := make(map[string]struct{})
	ids := []string{}

	for _, id := range ss {
		if _, ok := keys[id]; !ok {
			keys[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids
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
