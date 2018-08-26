package main

import "sync"

type Subscriber interface {
	OnMessage(msg string)
}

type PubSub struct {
	ss sync.Map
}

func (ps *PubSub) Publish(msg string) {
	ps.ss.Range(func(k, v interface{}) bool {
		s := v.(Subscriber)
		s.OnMessage(msg)
		return true
	})
}

func (ps *PubSub) Subscribe(id string, s Subscriber) {
	ps.ss.Store(id, s)
}
