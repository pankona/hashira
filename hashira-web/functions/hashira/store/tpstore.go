package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pankona/hashira/hashira-web/functions/hashira"
)

type TaskAndPriorityStore struct{}

func NewTaskAndPriorityStore() *TaskAndPriorityStore {
	return &TaskAndPriorityStore{}
}

func (t *TaskAndPriorityStore) Save(ctx context.Context, uid string, tp hashira.TaskAndPriority) error {
	client, err := firestore.NewClient(ctx, "hashira-web")
	if err != nil {
		return fmt.Errorf("failed to create firebase client: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(tp); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	if _, err := client.Collection("tasksAndPriorities").Doc(uid).Set(ctx, tp); err != nil {
		return fmt.Errorf("failed to write documents: %w", err)
	}

	return nil
}
