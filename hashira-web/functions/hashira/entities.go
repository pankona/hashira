package hashira

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type TaskAndPriority struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}
