package main

import (
	"context"

	"github.com/pankona/hashira/service"
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
	defer c.Update(context.Background())

	var err error

	switch event {
	// TODO: event should be dispatched by type assertion?
	case "add":
		err = c.m.hashirac.Create(context.Background(), data.(*service.Task))
	case "update":
		err = c.m.hashirac.Update(context.Background(), data.(*service.Task))
	case "delete":
		t := data.(*service.Task)
		err = c.m.hashirac.Delete(context.Background(), t.Id)
	case "updatePriority":
		p := data.([]*service.Priority)
		_, err = c.m.hashirac.UpdatePriority(context.Background(), p)
	default:
		// nop
	}

	return err
}

func (c *Ctrl) Update(ctx context.Context) error {
	tasks, err := c.m.List(ctx)
	if err != nil {
		return err
	}

	priorities, err := c.m.RetrievePriority(ctx)
	if err != nil {
		return err
	}

	c.pub.Publish("update", tasks, priorities)
	return nil
}
