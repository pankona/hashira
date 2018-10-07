package main

// Delegater in an interface to call delegate function,
// to cover functionality that is not covered by view
type Delegater interface {
	Delegate(event delegateEvent, data ...interface{}) error
}

type delegateEvent int

const (
	Add delegateEvent = iota
	Update
	Delete
	UpdatePriority
	UpdateBulk
)
