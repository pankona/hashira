package main

// Delegater in an interface to call delegate function,
// to cover functionality that is not covered by view
type Delegater interface {
	Delegate(event delegateEvent, data ...interface{}) error
}

type delegateEvent int

const (
	AddTask delegateEvent = iota
	UpdateTask
	DeleteTask
	UpdatePriority
	UpdateBulk
)
