package main

import "sync"

// Subscriber is an interface of PubSub subscriber
type Subscriber interface {
	OnEvent(event string, data ...interface{})
}

// Publisher is an interface of PubSub publisher
type Publisher interface {
	Publish(event string, data ...interface{})
}

// PubSub is a struct for pubsub message passing
type PubSub struct {
	ss sync.Map
}

// Publish publishes specified event and data to subscribers
func (ps *PubSub) Publish(event string, data ...interface{}) {
	ps.ss.Range(func(k, v interface{}) bool {
		s := v.(Subscriber)
		s.OnEvent(event, data...)
		return true
	})
}

// Subscribe registers subscriber to PubSub
func (ps *PubSub) Subscribe(id string, s Subscriber) {
	ps.ss.Store(id, s)
}
