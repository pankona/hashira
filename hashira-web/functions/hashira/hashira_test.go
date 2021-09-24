package hashira

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
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
