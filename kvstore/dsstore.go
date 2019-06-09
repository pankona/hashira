package kvstore

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/datastore"
)

// TODO: this should not be placed here
const dsProjectName = "hashira-auth"

type DSStore struct{}

type entity struct {
	value []byte
}

func (s *DSStore) Store(bucket, k string, v interface{}) {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, dsProjectName)
	if err != nil {
		// Handle error.
	}

	key := datastore.NameKey(bucket, k, nil)
	buf, err := json.Marshal(v)
	if err != nil {
		// TODO: error handling
		panic(err)
	}

	e := &entity{value: buf}
	if _, err := dsClient.Put(ctx, key, e); err != nil {
		// TODO: error handling
		panic(err)
	}
}

func (s *DSStore) Load(bucket, k string) (interface{}, bool) {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, dsProjectName)
	if err != nil {
		// TODO: error handling
		panic(err)
	}

	key := datastore.NameKey(bucket, k, nil)
	e := entity{}
	if err := dsClient.Get(ctx, key, &e); err != nil {
		// TODO: error handling
		return nil, false
	}

	var v interface{}
	if err := json.Unmarshal(e.value, &v); err != nil {
		// TODO: error handling
		panic(err)
	}

	return v, true
}
