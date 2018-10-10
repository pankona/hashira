package main

// Delegater in an interface to call delegate function,
// to cover functionality that is not covered by view
type Delegater interface {
	Delegate(event delegateEvent, data ...interface{}) error
}

type delegateEvent int

const (
	// AddTask is a event for adding a new task
	AddTask delegateEvent = iota
	// UpdateTask is a event for updating a task
	UpdateTask
	// DeleteTask is a event for deleting a task
	DeleteTask
	// UpdatePriority is a event for updating task priority
	UpdatePriority
	// UpdateBulk updates specified task and priority
	UpdateBulk
)
