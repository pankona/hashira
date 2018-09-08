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

func (c *Ctrl) Delegate(event string, data interface{}) error {
	switch event {
	case "add":
		err := c.m.hashirac.Create(context.Background(), data.(string))
		if err != nil {
			return err
		}
		return c.Update(context.Background())
	default:
	}
	return nil
}

func (c *Ctrl) Update(ctx context.Context) error {
	tasks, err := c.m.List(ctx)
	if err != nil {
		return err
	}

	c.pub.Publish("update", tasks)

	return nil
}
