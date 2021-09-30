package hashira

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	var (
		defaultBody = TaskAndPriority{
			Tasks: map[string]Task{
				"task1": {ID: "task1", Name: "task", Place: "Backlog", IsDeleted: false},
			},
			Priority: map[string][]string{
				"Backlog": {"task1"},
			},
		}
		defaultMockFunc = func(ctx context.Context, uid string) (TaskAndPriority, error) { return defaultBody, nil }
	)

	tests := []struct {
		name         string
		wantStatus   int
		wantRespBody TaskAndPriority
		mockFunc     func(ctx context.Context, uid string) (TaskAndPriority, error)
	}{
		{
			name:         "regular case",
			wantStatus:   http.StatusOK,
			wantRespBody: defaultBody,
			mockFunc:     defaultMockFunc,
		},
		{
			name:         "failed to load",
			wantStatus:   http.StatusInternalServerError,
			wantRespBody: defaultBody,
			mockFunc: func(ctx context.Context, uid string) (TaskAndPriority, error) {
				return TaskAndPriority{}, errors.New("dummy error")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add("Authorization", "bearer 123")

			h := &Hashira{
				AccessTokenStore: &mockAccessTokenStore{
					mockFindUidByAccessToken: func(ctx context.Context, accesstoken string) (string, error) {
						return "dummy_uid", nil
					},
				},
				TaskAndPriorityStore: &mockTaskAndPriorityStore{
					mockLoad: tt.mockFunc,
				},
			}

			h.Download(rec, req)

			result := rec.Result()
			if result.StatusCode != tt.wantStatus {
				t.Errorf("unexpected result: [got] %d [want] %d", result.StatusCode, tt.wantStatus)
			}

			// don't continue testing since body is nil if response is not ok
			if tt.wantStatus != http.StatusOK {
				return
			}

			var gotRespBody TaskAndPriority
			if err := json.NewDecoder(result.Body).Decode(&gotRespBody); err != nil {
				t.Errorf("failed to decode response body: %v", err)
			}

			if diff := cmp.Diff(tt.wantRespBody, gotRespBody); diff != "" {
				t.Errorf("unexpected result: diff: %s", diff)
			}
		})
	}
}
