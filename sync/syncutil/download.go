package syncutil

import (
	"context"
	"fmt"
	"strconv"

	hc "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
	"github.com/pankona/hashira/sync"
)

func (c *Client) Download(accesstoken string) error {
	sc := sync.NewClient()
	result, err := sc.Download(accesstoken)
	if err != nil {
		return fmt.Errorf("failed to download task and priority: %w", err)
	}

	cli := &hc.Client{Address: "localhost:" + strconv.Itoa(c.DaemonPort)}
	for _, task := range result.Tasks {
		err = cli.Update(context.Background(), &service.Task{
			Id:        task.ID,
			Name:      task.Name,
			Place:     service.Place(service.Place_value[task.Place]),
			IsDeleted: task.IsDeleted,
			IsDirty:   false,
		})
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}
	}

	newPriorities := map[string]*service.Priority{}
	for k, v := range result.Priority {
		newPriorities[k] = &service.Priority{Ids: v, IsDirty: false}
	}

	oldPriorities, err := cli.RetrievePriority(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve old priorities: %w", err)
	}

	priorities := mergePriorities(newPriorities, oldPriorities)

	_, err = cli.UpdatePriority(context.Background(), priorities)
	if err != nil {
		return fmt.Errorf("failed to update priority: %w", err)
	}

	return nil
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
