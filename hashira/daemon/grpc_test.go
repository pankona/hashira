package daemon

import (
	"context"
	"testing"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
)

type mockDB struct {
	database.Databaser
	data map[string][]byte
}

func (m *mockDB) Save(id string, value []byte) error {
	m.data[id] = value
	return nil
}

func TestEndPointCreate(t *testing.T) {
	d := &Daemon{
		DB: &mockDB{data: make(map[string][]byte)},
	}

	cc := &service.CommandCreate{
		Name:  "test",
		Place: service.Place_BACKLOG,
	}
	result, err := d.Create(context.Background(), cc)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %s", err.Error())
	}

	if result.GetTask().GetName() != "test" {
		t.Errorf("unexpected result. [want] %v [got] %v", result.GetTask().GetName(), "test")
	}

	if result.GetTask().GetPlace() != service.Place_BACKLOG {
		t.Errorf("unexpected result. [want] %v [got] %v", result.GetTask().GetPlace(), service.Place_BACKLOG)
	}
}

func TestEndPointRetrieve(t *testing.T) {
}
