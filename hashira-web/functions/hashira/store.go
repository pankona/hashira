package hashira

import "context"

type AccessTokenStore interface {
	FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error)
}

type TaskAndPriorityStore interface {
	Save(ctx context.Context, uid string, tp TaskAndPriority) error
	Load(ctx context.Context, uid string) (TaskAndPriority, error)
}
