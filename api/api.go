package api

type TaskStore interface {
	SaveTasks(userID string, tasks Tasks) error
	SavePriority(userID string, priority Priority) error
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

type Tasks []Task

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
