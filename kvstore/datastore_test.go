package kvstore

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"cloud.google.com/go/datastore"
)

type testEntity struct {
	Value []byte
}

// Launch DataStore emulator in advance to run this test.
// DATASTORE_EMULATOR_HOST is not set as a environment variable,
// this function forcibly goes failure.
// $ export DATASTORE_EMULATOR_HOST=localhost:8081
func TestHowToUseDataStore(t *testing.T) {
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		t.Fatalf("Run DataStore emulator and configure environment variable of \"DATASTORE_EMULATOR_HOST\" in advance to run this test.")
	}

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, "hashira-auth")
	if err != nil {
		t.Fatalf("failed to create datastore client: %v", err)
	}

	key := datastore.NameKey("testkind", "testkey", nil)
	testStr := "hogehoge"

	buf, err := json.Marshal(testStr)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	e := &testEntity{Value: buf}
	if _, err := dsClient.Put(ctx, key, e); err != nil {
		// TODO: error handling
		t.Fatalf("failed to put: %v", err)
	}

	e2 := testEntity{}
	if err := dsClient.Get(ctx, key, &e2); err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	var ret interface{}
	err = json.Unmarshal(e2.Value, &ret)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if ret != testStr {
		t.Fatalf("unexpected value returned from Get: %v, expected: %v", ret, testStr)
	}
}

func TestStoreAndLoad(t *testing.T) {
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
