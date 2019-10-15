package store

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/datastore"
)

type testEntity struct {
	Value string
}

// Launch DataStore emulator in advance to run this test.
// DATASTORE_EMULATOR_HOST is not set as a environment variable,
// this function forcibly goes failure.
// $ export DATASTORE_EMULATOR_HOST=localhost:8081
func testHowToUseDataStore(t *testing.T) {
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		t.Fatalf("Run DataStore emulator and configure environment variable " +
			"\"DATASTORE_EMULATOR_HOST\" in advance to run this test.")
	}

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, "hashira-auth")
	if err != nil {
		t.Fatalf("failed to create datastore client: %v", err)
	}

	key := datastore.NameKey("testkind", "testkey", nil)
	e1 := &testEntity{Value: "hogehoge"}

	// Thus dsClient.Put takes interface{}, pointer types can be accepted
	if _, err := dsClient.Put(ctx, key, e1); err != nil {
		// TODO: error handling
		t.Fatalf("failed to put: %v", err)
	}

	e2 := testEntity{}
	if err := dsClient.Get(ctx, key, &e2); err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if e1.Value != e2.Value {
		t.Fatalf("unexpected value returned from Get: %v, expected: %v", e1.Value, e2.Value)
	}
}

func testStoreAndLoad(t *testing.T) {
	ds := &DSStore{}
	ds.Store("bucket", "key", "value")
	v, ok := ds.Load("bucket", "key")
	if !ok {
		t.Fatalf("failed to load. fatal.")
	}
	if "value" != v.(string) {
		t.Fatalf("unexpected result. [want] %s [got] %s", "value", v.(string))
	}
}

func TestDataStoreUsage(t *testing.T) {
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		t.Logf("Run DataStore emulator and configure environment variable " +
			"\"DATASTORE_EMULATOR_HOST\" in advance to run this test.")
		return
	}

	testHowToUseDataStore(t)
	testStoreAndLoad(t)
}
