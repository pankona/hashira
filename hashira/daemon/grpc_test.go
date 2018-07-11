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

func (m *mockDB) ForEach(f func(k, v []byte) error) error {
	for k, v := range m.data {
		err := f([]byte(k), v)
		if err != nil {
			return err
		}
	}
	return nil
}

var tcs = []struct {
	inName    string
	inPlace   service.Place
	wantName  string
	wantPlace service.Place
}{
	{
		inName: "test", inPlace: service.Place_BACKLOG,
		wantName: "test", wantPlace: service.Place_BACKLOG,
	},
}

func testCreate(d *Daemon, t *testing.T) {
	for _, tc := range tcs {
		cc := &service.CommandCreate{
			Name:  tc.inName,
			Place: tc.inPlace,
		}
		result, err := d.Create(context.Background(), cc)
		if err != nil {
			t.Fatalf("Create returned unexpected error: %s", err.Error())
		}

		if result.GetTask().GetName() != tc.wantName {
			t.Errorf("unexpected result. [got] %v [want] %v", result.GetTask().GetName(), tc.wantName)
		}

		if result.GetTask().GetPlace() != tc.wantPlace {
			t.Errorf("unexpected result. [got] %v [want] %v", result.GetTask().GetPlace(), tc.wantPlace)
		}
	}
}

func TestEndPointCreate(t *testing.T) {
	d := &Daemon{
		DB: &mockDB{data: make(map[string][]byte)},
	}
	testCreate(d, t)
}

func TestEndPointRetrieve(t *testing.T) {
	d := &Daemon{
		DB: &mockDB{data: make(map[string][]byte)},
	}
	testCreate(d, t)

	rc := &service.CommandRetrieve{}
	result, err := d.Retrieve(context.Background(), rc)
	if err != nil {
		t.Fatalf("Retrieve returned unexpected error: %s", err.Error())
	}

	tasks := result.GetTasks()
	if len(tasks) != len(tcs) {
		t.Errorf("unexpected result. [got] %d [want] %d", len(tasks), len(tcs))
	}
}
