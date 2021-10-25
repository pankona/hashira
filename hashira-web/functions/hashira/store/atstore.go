package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

type AccessTokenStore struct{}

func NewAccessTokenStore() *AccessTokenStore {
	return &AccessTokenStore{}
}

func (a *AccessTokenStore) FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error) {
	// accesstoken length is assumed to be 36 (uuid).
	// If the accesstoken is longer than it, let's assume it is idtoken
	if len(accesstoken) > 36 {
		app, err := firebase.NewApp(ctx, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create new firebase app: %w", err)
		}

		cli, err := app.Auth(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to prepare Auth client: %w", err)
		}

		token, err := cli.VerifyIDToken(ctx, accesstoken)
		if err != nil {
			return "", fmt.Errorf("failed to verify idtoken: %w", err)
		}

		return token.UID, nil
	}

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
