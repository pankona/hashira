package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

// Entity is a sample entity
type Entity struct {
	// Value is a sample value
	Value string
	small string
}

func main() {
	ctx := context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	dsClient, err := datastore.NewClient(ctx, "my-project")
	if err != nil {
		panic(err)
	}

	k := datastore.NameKey("Entity", "stringID", nil)
	e := new(Entity)
	if err := dsClient.Get(ctx, k, e); err != nil {
		panic(err)
	}

	old := e.Value
	e.Value = "Hello World!"
	e.small = "small"

	if _, err := dsClient.Put(ctx, k, e); err != nil {
		panic(err)
	}

	fmt.Printf("Updated value from %q to %q\n", old, e)
}
