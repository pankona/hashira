package api

type TaskStore interface {
	SaveTasks(userID string, tasks Tasks) error
	SavePriority(userID string, priority Priority) error
	LoadTasks(userID string) (Tasks, error)
	LoadPriority(userID string) (Priority, error)
}

type Place string

const (
	Backlog Place = "BACKLOG"
	ToDo    Place = "TODO"
	Doing   Place = "DOING"
	Done    Place = "DONE"
)

type Task struct {
	ID        string
	Name      string
	Place     Place
	IsDeleted bool
}

type Tasks map[string]Task

type Priority map[Place][]string

type API struct {
	TaskStore TaskStore
}

func (a *API) Upload(userID string, tasks Tasks, p Priority) error {
	err := a.TaskStore.SaveTasks(userID, tasks)
	if err != nil {
		return err
	}
	return a.TaskStore.SavePriority(userID, p)
}

func (a *API) Download(userID string) (Tasks, Priority, error) {
	tasks, err := a.TaskStore.LoadTasks(userID)
	if err != nil {
		return Tasks{}, Priority{}, err
	}
	priority, err := a.TaskStore.LoadPriority(userID)
	if err != nil {
		return Tasks{}, Priority{}, err
	}
	return tasks, priority, nil
}
