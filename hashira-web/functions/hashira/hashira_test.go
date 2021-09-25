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

type mockAccessTokenStore struct {
	mockFindUidByAccessToken func(ctx context.Context, accesstoken string) (string, error)
}

func (m *mockAccessTokenStore) FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error) {
	return m.mockFindUidByAccessToken(ctx, accesstoken)
}

func TestTestAccessToken(t *testing.T) {
	t.Parallel()

	defaultMockFindByAccessToken := func(ctx context.Context, accesstoken string) (string, error) {
		return "dummy_uid", nil
	}

	tests := []struct {
		name       string
		inHeaders  []string
		wantStatus int
		mockFunc   func(ctx context.Context, accesstoken string) (string, error)
	}{
		{
			name:       "regular case",
			inHeaders:  []string{"bearer 123"},
			mockFunc:   defaultMockFindByAccessToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "header missing",
			inHeaders:  nil,
			mockFunc:   defaultMockFindByAccessToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "too many headers",
			inHeaders:  []string{"bearer 123", "bearer 456", "bearer 789"},
			mockFunc:   defaultMockFindByAccessToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no bearer",
			inHeaders:  []string{"123"},
			mockFunc:   defaultMockFindByAccessToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "associated user not found ",
			inHeaders: []string{"bearer nobody-use-this-access-token"},
			mockFunc: func(ctx context.Context, accesstoken string) (string, error) {
				return "", errors.New("user not found who uses this accesstoken")
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			for _, header := range tt.inHeaders {
				req.Header.Add("Authorization", header)
			}

			h := &Hashira{
				AccessTokenStore: &mockAccessTokenStore{
					mockFindUidByAccessToken: tt.mockFunc,
				},
			}

			h.TestAccessToken(rec, req)

			result := rec.Result()
			if result.StatusCode != tt.wantStatus {
				t.Errorf("unexpected result: [got] %d [want] %d", result.StatusCode, tt.wantStatus)
			}
		})
	}
}

type mockTaskAndPriorityStore struct {
	mockSave func(ctx context.Context, uid string, tp TaskAndPriority) error
	mockLoad func(ctx context.Context, uid string) (TaskAndPriority, error)
}

func (m *mockTaskAndPriorityStore) Save(ctx context.Context, uid string, tp TaskAndPriority) error {
	return m.mockSave(ctx, uid, tp)
}

func (m *mockTaskAndPriorityStore) Load(ctx context.Context, uid string) (TaskAndPriority, error) {
	return m.mockLoad(ctx, uid)
}

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
