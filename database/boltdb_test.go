package database

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func setup(filename *string) (func(), error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, errors.New("fatal. failed to create tempfile")
	}
	*filename = f.Name()
	err = f.Close()
	if err != nil {
		return nil, errors.New("fatal. failed to close file")
	}
	return func() {
		err = os.RemoveAll(*filename)
		if err != nil {
			fmt.Printf("os.RemoveAll returned error: %v", err)
		}
	}, nil
}

func TestBoltDBInterface(t *testing.T) {
	var _ Databaser = &BoltDB{}
}

func TestBoltDBInitialize(t *testing.T) {
	db := &BoltDB{}
	var filename string
	teardown, err := setup(&filename)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer teardown()

	err = db.Initialize(filename)
	if err != nil {
		t.Fatalf("Initialized returned unexpected error: %s", err.Error())
	}
}

func TestBoltDBFinalize(t *testing.T) {
	db := &BoltDB{}
	err := db.Finalize()
	if err != nil {
		t.Fatalf("Finalize returned unexpected error: %s", err.Error())
	}
}

func TestBoltDBSaveLoad(t *testing.T) {
	db := &BoltDB{}
	var filename string
	teardown, err := setup(&filename)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer teardown()

	err = db.Initialize(filename)
	if err != nil {
		t.Fatalf("Initialized returned unexpected error: %s", err.Error())
	}

	_, err = db.Save("testbucket", "testid", []byte("testdata"))
	if err != nil {
		t.Fatalf("save returned unexpected error: %s", err.Error())
	}

	v, err := db.Load("testbucket", "testid")
	if err != nil {
		t.Fatalf("load returned unexpected error: %s", err.Error())
	}

	if string(v) != "testdata" {
		t.Fatalf("load returned unexpected result. [got] %s [want] %s", string(v), "testdata")
	}
}

func TestBoltDBSaveLoadWithoutID(t *testing.T) {
	db := &BoltDB{}
	var filename string
	teardown, err := setup(&filename)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer teardown()

	err = db.Initialize(filename)
	if err != nil {
		t.Fatalf("Initialized returned unexpected error: %s", err.Error())
	}

	for i := 0; i < 10; i++ {
		_, err = db.Save("testbucket", "", []byte("testdata"))
		if err != nil {
			t.Fatalf("save returned unexpected error: %s", err.Error())
		}
	}

	err = db.ForEach("testbucket", func(k, v []byte) error {
		t.Logf("[%s] %s", string(k), string(v))
		return nil
	})
	if err != nil {
		t.Fatalf("ForEach returned unexpected error: %s", err.Error())
	}
}
