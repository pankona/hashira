package kvstore

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/datastore"
)

// TODO: this should not be placed here
const dsProjectName = "hashira-auth"

// DSStore is an struct that implements KVStore interface using google cloud datastore
type DSStore struct{}

type entity struct {
	Value []byte
}

// Store stores value associated with key to bucket
func (s *DSStore) Store(bucket, k string, v interface{}) {
	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, dsProjectName)
	if err != nil {
		// Handle error.
		panic(err)
	}

	key := datastore.NameKey(bucket, k, nil)
	buf, err := json.Marshal(v)
	if err != nil {
		// TODO: error handling
		panic(err)
	}

	e := &entity{Value: buf}
	if _, err := dsClient.Put(ctx, key, e); err != nil {
		// TODO: error handling
		panic(err)
	}
}

// Load loads key's value from bucket
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
	if err := json.Unmarshal(e.Value, &v); err != nil {
		// TODO: error handling
		return nil, false
	}

	return v, true
}
