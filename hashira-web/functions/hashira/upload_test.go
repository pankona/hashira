package hashira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	var (
		defaultMockFunc = func(ctx context.Context, uid string, tp TaskAndPriority) error { return nil }
		defaultBody     = TaskAndPriority{
			Tasks: map[string]Task{
				"task1": {ID: "task1", Name: "task", Place: "Backlog", IsDeleted: false},
			},
			Priority: map[string][]string{
				"Backlog": {"task1"},
			},
		}
	)

	tests := []struct {
		name       string
		inBody     TaskAndPriority
		wantStatus int
		mockFunc   func(ctx context.Context, uid string, tp TaskAndPriority) error
	}{
		{
			name:       "regular case",
			inBody:     defaultBody,
			wantStatus: http.StatusOK,
			mockFunc:   defaultMockFunc,
		},
		{
			name:       "failed to save",
			inBody:     defaultBody,
			wantStatus: http.StatusInternalServerError,
			mockFunc:   func(ctx context.Context, uid string, tp TaskAndPriority) error { return errors.New("dummy error") },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()

			buf := &bytes.Buffer{}
			if err := json.NewEncoder(buf).Encode(tt.inBody); err != nil {
				t.Fatalf("failed to prepare request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/", buf)
			req.Header.Add("Authorization", "bearer 123")

			h := &Hashira{
				AccessTokenStore: &mockAccessTokenStore{
					mockFindUidByAccessToken: func(ctx context.Context, accesstoken string) (string, error) {
						return "dummy_uid", nil
					},
				},
				TaskAndPriorityStore: &mockTaskAndPriorityStore{
					mockSave: tt.mockFunc,
				},
			}

			h.Upload(rec, req)

			result := rec.Result()
			if result.StatusCode != tt.wantStatus {
				t.Errorf("unexpected result: [got] %d [want] %d", result.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestMergeTaskAndPriorities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		inNewTP TaskAndPriority
		inOldTP TaskAndPriority
		wantTP  TaskAndPriority
	}{
		{
			name: "a",
			inNewTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task1": {ID: "task1", Name: "task1", Place: "BACKLOG", IsDeleted: false},
				},
				Priority: map[string][]string{
					"BACKLOG": {"task1"},
					"TODO":    {},
					"DOING":   {},
					"DONE":    {},
				},
			},
			inOldTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task2": {ID: "task2", Name: "task2", Place: "BACKLOG", IsDeleted: false},
				},
				Priority: map[string][]string{
					"BACKLOG": {"task2"},
				},
			},
			wantTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task1": {ID: "task1", Name: "task1", Place: "BACKLOG", IsDeleted: false},
					"task2": {ID: "task2", Name: "task2", Place: "BACKLOG", IsDeleted: false},
				},
				Priority: map[string][]string{
					"BACKLOG": {"task1", "task2"},
					"TODO":    {},
					"DOING":   {},
					"DONE":    {},
				},
			},
		},
		{
			name: "b",
			inNewTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task1": {ID: "task1", Name: "task1", Place: "TODO", IsDeleted: false},
				},
				Priority: map[string][]string{
					"TODO": {"task1"},
				},
			},
			inOldTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task1": {ID: "task1", Name: "task1", Place: "BACKLOG", IsDeleted: false},
				},
				Priority: map[string][]string{
					"BACKLOG": {"task1"},
				},
			},
			wantTP: TaskAndPriority{
				Tasks: map[string]Task{
					"task1": {ID: "task1", Name: "task1", Place: "TODO", IsDeleted: false},
				},
				Priority: map[string][]string{
					"BACKLOG": {},
					"TODO":    {"task1"},
					"DOING":   {},
					"DONE":    {},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTP, err := mergeTaskAndPriorities(tt.inNewTP, tt.inOldTP)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(gotTP, tt.wantTP); diff != "" {
				t.Errorf("unexpected result: diff: %v", diff)
			}
		})
	}
}

func TestMergePriorities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		inNewP map[string][]string
		inOldP map[string][]string
		wantP  map[string][]string
	}{
		{
			name: "a",
			inNewP: map[string][]string{
				"TODO": {
					"task1",
				},
			},
			inOldP: map[string][]string{
				"TODO": {
					"task2",
				},
			},
			wantP: map[string][]string{
				"BACKLOG": {},
				"TODO": {
					"task1",
					"task2",
				},
				"DOING": {},
				"DONE":  {},
			},
		},
		{
			name: "b",
			inNewP: map[string][]string{
				"BACKLOG": {
					"task1",
				},
			},
			inOldP: map[string][]string{},
			wantP: map[string][]string{
				"BACKLOG": {
					"task1",
				},
				"TODO":  {},
				"DOING": {},
				"DONE":  {},
			},
		},
		{
			name: "c",
			inNewP: map[string][]string{
				"TODO": {
					"task2",
				},
			},
			inOldP: map[string][]string{
				"BACKLOG": {
					"task1",
				},
			},
			wantP: map[string][]string{
				"BACKLOG": {
					"task1",
				},
				"TODO": {
					"task2",
				},
				"DOING": {},
				"DONE":  {},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotP := mergePriorities(tt.inNewP, tt.inOldP)
			if diff := cmp.Diff(gotP, tt.wantP); diff != "" {
				t.Errorf("unexpected result: diff: %v", diff)
			}
		})
	}
}
