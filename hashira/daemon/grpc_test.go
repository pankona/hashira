package daemon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
)

type mockDB struct {
	database.Databaser
	data map[string]map[string][]byte
}

func (m *mockDB) Save(bucket, id string, value []byte) (string, error) {
	if m.data == nil {
		m.data = make(map[string]map[string][]byte)
	}

	if m.data[bucket] == nil {
		m.data[bucket] = make(map[string][]byte)
	}

	m.data[bucket][id] = value
	return id, nil
}

func (m *mockDB) Load(bucket, id string) ([]byte, error) {
	if m.data == nil {
		return nil, nil
	}

	if m.data[bucket] == nil {
		return nil, nil
	}

	return m.data[bucket][id], nil
}

func (m *mockDB) ForEach(bucket string, f func(k, v []byte) error) error {
	for k, v := range m.data[bucket] {
		err := f([]byte(k), v)
		if err != nil {
			return err
		}
	}
	return nil
}

var tcs = []struct {
	inName    string
	wantName  string
	inPlace   service.Place
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
			Task: &service.Task{
				Name:  tc.inName,
				Place: tc.inPlace,
			},
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

func TestUpdatePriority(t *testing.T) {
	p := map[string]*service.Priority{
		service.Place_BACKLOG.String(): {
			Ids: []string{"1", "2", "3", "4", "5"},
		},
		service.Place_TODO.String(): {
			Ids: []string{"6", "7", "8", "9", "10"},
		},
		service.Place_DOING.String(): {
			Ids: []string{"11", "12", "13", "14", "15"},
		},
		service.Place_DONE.String(): {
			Ids: []string{"16", "17", "18", "19", "20"},
		},
	}

	d := &Daemon{
		DB: &mockDB{},
	}
	m, err := d.updatePriority(p)
	require.Nil(t, err)

	require.Equal(t, m[service.Place_BACKLOG.String()].Ids, []string{"1", "2", "3", "4", "5"})
	require.Equal(t, m[service.Place_TODO.String()].Ids, []string{"6", "7", "8", "9", "10"})
	require.Equal(t, m[service.Place_DOING.String()].Ids, []string{"11", "12", "13", "14", "15"})
	require.Equal(t, m[service.Place_DONE.String()].Ids, []string{"16", "17", "18", "19", "20"})
}

func TestEndPointCreate(t *testing.T) {
	d := &Daemon{
		DB: &mockDB{},
	}
	testCreate(d, t)
}

func TestEndPointRetrieve(t *testing.T) {
	d := &Daemon{
		DB: &mockDB{},
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
