package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pankona/hashira/hashira-web/functions/hashira"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (t *TaskAndPriorityStore) Load(ctx context.Context, uid string) (hashira.TaskAndPriority, error) {
	client, err := firestore.NewClient(ctx, "hashira-web")
	if err != nil {
		return hashira.TaskAndPriority{}, fmt.Errorf("failed to create firebase client: %w", err)
	}

	ds, err := client.Collection("tasksAndPriorities").Doc(uid).Get(ctx)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.NotFound {
				// update with empty tasks and priorities for initial data
				tp := hashira.TaskAndPriority{}
				buf := &bytes.Buffer{}
				if err := json.NewEncoder(buf).Encode(tp); err != nil {
					return hashira.TaskAndPriority{}, fmt.Errorf("failed to encode data: %w", err)
				}
				if _, err := client.Collection("tasksAndPriorities").Doc(uid).Create(ctx, tp); err != nil {
					return hashira.TaskAndPriority{}, fmt.Errorf("failed to create initial documents: %w", err)
				}
				return hashira.TaskAndPriority{}, nil
			}
		}
		return hashira.TaskAndPriority{}, fmt.Errorf("failed to get documents: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(ds.Data()); err != nil {
		return hashira.TaskAndPriority{}, fmt.Errorf("failed to encode data: %w", err)
	}

	var ret hashira.TaskAndPriority
	json.NewDecoder(buf).Decode(&ret)

	return ret, nil
}
