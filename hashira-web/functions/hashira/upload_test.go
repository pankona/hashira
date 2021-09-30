package hashira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
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
