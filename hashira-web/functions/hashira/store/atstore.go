package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type AccessTokenStore struct{}

func NewAccessTokenStore() *AccessTokenStore {
	return &AccessTokenStore{}
}

func (a *AccessTokenStore) FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error) {
	client, err := firestore.NewClient(ctx, "hashira-web")
	if err != nil {
		return "", fmt.Errorf("failed to create firebase client: %w", err)
	}

	iter := client.Collection("accesstokens").Where("accesstoken", "==", accesstoken).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to get documents: %w", err)
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("user not found who has the accesstoken [%v]", accesstoken)
	}
	if len(docs) >= 2 {
		return "", fmt.Errorf("at least two users has same accesstoken [%s]. data inconsistency. fatal", accesstoken)
	}
	doc := docs[0]
	uid, ok := doc.Data()["uid"]
	if !ok {
		return "", fmt.Errorf("fetched accesstoken data doesn't have uid field. fatal")
	}

	ret, ok := uid.(string)
	if !ok {
		return "", fmt.Errorf("fetched uid is not string. fatal")
	}

	return ret, nil
}
