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

	tcs := []struct {
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

func TestEndPointRetrieve(t *testing.T) {
}
