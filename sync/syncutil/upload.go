package syncutil

import (
	"context"
	"fmt"
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

func (c *Client) Upload(accesstoken string, uploadTarget UploadTarget) error {
	cli := &hc.Client{Address: "localhost:" + strconv.Itoa(c.DaemonPort)}
	allTasks, err := cli.RetrieveAll(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	allPriorities, err := cli.RetrievePriority(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve priorities: %w", err)
	}

	ur := newUploadRequest(allTasks, allPriorities, uploadTarget)

	if len(ur.Tasks) == 0 && !isPriorityDirty(allPriorities) {
		// there's no task to upload
		return fmt.Errorf("there's no dirty task. no task to upload")
	}

	sc := sync.NewClient()
	err = sc.Upload(accesstoken, ur)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}

	for _, task := range allTasks {
		if !task.IsDeleted {
			continue
		}
		if err := cli.PhysicalDelete(context.Background(), task.Id); err != nil {
			return fmt.Errorf("failed to physical delete a task: %w", err)
		}
	}

	return nil
}

func isPriorityDirty(p map[string]*service.Priority) bool {
	for _, v := range p {
		if v.IsDirty {
			return true
		}
	}
	return false
}
