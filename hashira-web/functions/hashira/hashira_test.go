package hashira

import (
	"context"
)

type mockAccessTokenStore struct {
	mockFindUidByAccessToken func(ctx context.Context, accesstoken string) (string, error)
}

func (m *mockAccessTokenStore) FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error) {
	return m.mockFindUidByAccessToken(ctx, accesstoken)
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
