package syncutil

import (
	"context"
	"log"
	"strconv"

	hc "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
	"github.com/pankona/hashira/sync"
)

type UploadTarget int

const (
	UploadAll UploadTarget = iota
	UploadDirtyOnly
)

func newUploadRequest(tasks map[string]*service.Task, priorities map[string]*service.Priority, uploadTarget UploadTarget) sync.UploadRequest {
	ur := sync.UploadRequest{
		Tasks: map[string]sync.Task{},
	}

	for k, v := range tasks {
		if uploadTarget == UploadDirtyOnly && !v.IsDirty {
			continue
		}
		ur.Tasks[k] = sync.Task{
			ID:        v.Id,
			Name:      v.Name,
			Place:     v.Place.String(),
			IsDeleted: v.IsDeleted,
		}
	}

	ur.Priority = sync.Priority{}
	for k, v := range priorities {
		ur.Priority[k] = v.Ids
	}

	return ur
}

func (c *Client) Upload(accesstoken string, uploadTarget UploadTarget) {
	log.Println("upload started")

	cli := &hc.Client{Address: "localhost:" + strconv.Itoa(c.DaemonPort)}
	allTasks, err := cli.RetrieveAll(context.Background())
	if err != nil {
		log.Printf("failed to retrieve tasks: %v", err)
		return
	}
	allPriorities, err := cli.RetrievePriority(context.Background())
	if err != nil {
		log.Printf("failed to retrieve priorities: %v", err)
		return
	}

	ur := newUploadRequest(allTasks, allPriorities, uploadTarget)

	if len(ur.Tasks) == 0 && !isPriorityDirty(allPriorities) {
		// there's no task to upload
		log.Println("there's no dirty task. no task to upload")
		return
	}

	sc := sync.NewClient()
	err = sc.Upload(accesstoken, ur)
	if err != nil {
		log.Printf("failed to upload: %v", err)
		return
	}

	for _, task := range allTasks {
		if !task.IsDeleted {
			continue
		}
		if err := cli.PhysicalDelete(context.Background(), task.Id); err != nil {
			log.Printf("failed to physical delete a task: %v", err)
		}
		log.Printf("task [%s] is deleted physically", task.Id)
	}

	log.Println("upload completed")
}

func isPriorityDirty(p map[string]*service.Priority) bool {
	for _, v := range p {
		if v.IsDirty {
			return true
		}
	}
	return false
}
