package daemon

import (
	"context"
	"testing"

	"github.com/pankona/hashira/database"
	"github.com/pankona/hashira/service"
	"github.com/stretchr/testify/require"
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

type testcase struct {
	inName    string
	wantName  string
	inPlace   service.Place
	wantPlace service.Place
}

func testcases() []testcase {
	tcs := []testcase{
		{inName: "test", inPlace: service.Place_BACKLOG, wantName: "test", wantPlace: service.Place_BACKLOG},
	}
	return tcs
}

func testCreate(d *Daemon, t *testing.T) {
	tcs := testcases()
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

	tcs := testcases()
	tasks := result.GetTasks()
	if len(tasks) != len(tcs) {
		t.Errorf("unexpected result. [got] %d [want] %d", len(tasks), len(tcs))
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		id   string
		want []string
	}{
		{
			name: "remove middle element",
			ids:  []string{"a", "b", "c"},
			id:   "b",
			want: []string{"a", "c"},
		},
		{
			name: "remove first element",
			ids:  []string{"a", "b", "c"},
			id:   "a",
			want: []string{"b", "c"},
		},
		{
			name: "remove last element",
			ids:  []string{"a", "b", "c"},
			id:   "c",
			want: []string{"a", "b"},
		},
		{
			name: "remove nonexistent element",
			ids:  []string{"a", "b", "c"},
			id:   "z",
			want: []string{"a", "b", "c"},
		},
		{
			name: "remove from empty slice",
			ids:  []string{},
			id:   "a",
			want: []string{},
		},
		{
			name: "remove only element",
			ids:  []string{"a"},
			id:   "a",
			want: []string{},
		},
		{
			name: "remove all duplicates",
			ids:  []string{"a", "a", "b"},
			id:   "a",
			want: []string{"b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := remove(tt.ids, tt.id)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRepair(t *testing.T) {
	mkTasks := func(ids ...string) map[string]*service.Task {
		m := make(map[string]*service.Task)
		for _, id := range ids {
			m[id] = &service.Task{Id: id}
		}
		return m
	}
	tests := []struct {
		name  string
		dst   []string
		tasks map[string]*service.Task
		want  []string
	}{
		{
			name:  "dst matches tasks exactly",
			dst:   []string{"a", "b"},
			tasks: mkTasks("a", "b"),
			want:  []string{"a", "b"},
		},
		{
			name:  "dst shorter than tasks adds missing",
			dst:   []string{"a"},
			tasks: mkTasks("a", "b", "c"),
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "dst longer than tasks removes extras",
			dst:   []string{"a", "b", "c", "d"},
			tasks: mkTasks("a", "b"),
			want:  []string{"a", "b"},
		},
		{
			name:  "empty dst adds all tasks",
			dst:   []string{},
			tasks: mkTasks("a", "b"),
			want:  []string{"a", "b"},
		},
		{
			name:  "empty tasks removes all dst",
			dst:   []string{"a", "b"},
			tasks: mkTasks(),
			want:  []string{},
		},
		{
			name:  "dst has stale ids interspersed",
			dst:   []string{"x", "a", "y", "b"},
			tasks: mkTasks("a", "b"),
			want:  []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repair(tt.dst, tt.tasks)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestLookupTaskByID(t *testing.T) {
	tasks := map[string]map[string]*service.Task{
		"BACKLOG": {
			"1": &service.Task{Id: "1"},
			"2": &service.Task{Id: "2"},
		},
		"TODO": {
			"3": &service.Task{Id: "3"},
		},
	}
	tests := []struct {
		name string
		id   string
		want *service.Task
	}{
		{
			name: "found in first pane",
			id:   "1",
			want: &service.Task{Id: "1"},
		},
		{
			name: "found in second pane",
			id:   "3",
			want: &service.Task{Id: "3"},
		},
		{
			name: "not found",
			id:   "99",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookupTaskByID(tasks, tt.id)
			require.Equal(t, tt.want, got)
		})
	}
}
