package main

import "sync"

type Subscriber interface {
	OnEvent(event string, data ...interface{})
}

type Publisher interface {
	Publish(event string, data ...interface{})
}

type PubSub struct {
	ss sync.Map
}

func (ps *PubSub) Publish(event string, data ...interface{}) {
	ps.ss.Range(func(k, v interface{}) bool {
		s := v.(Subscriber)
		s.OnEvent(event, data...)
		return true
	})
}

func (ps *PubSub) Subscribe(id string, s Subscriber) {
	ps.ss.Store(id, s)
}
