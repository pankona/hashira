package main

import (
	"context"
)

type Ctrl struct {
	m   *Model
	pub Publisher
}

func (c *Ctrl) Initialize() {
}

func (c *Ctrl) SetPublisher(p Publisher) {
	c.pub = p
}

func (c *Ctrl) Update(ctx context.Context) error {
	tasks, err := c.m.List(ctx)
	if err != nil {
		return err
	}

	c.pub.Publish("update", tasks)
	return nil
}
