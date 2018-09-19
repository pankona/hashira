package main

import (
	"context"
	"log"

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
	switch event {
	// TODO: event should be dispatched by type assertion?
	case "add":
		err := c.m.hashirac.Create(context.Background(), data.(string))
		if err != nil {
			return err
		}
		return c.Update(context.Background())
	case "delete":
		t := data.(*service.Task)
		err := c.m.hashirac.Delete(context.Background(), t.Id)
		if err != nil {
			return err
		}
		return c.Update(context.Background())
	case "updatePriority":
		p := data.([]*service.Priority)
		_, err := c.m.hashirac.UpdatePriority(context.Background(), p)
		if err != nil {
			return err
		}
		return c.Update(context.Background())
	default:
		// nop
	}
	return nil
}

func (c *Ctrl) Update(ctx context.Context) error {
	tasks, err := c.m.List(ctx)
	if err != nil {
		return err
	}

	priorities, err := c.m.RetrievePriority(ctx)
	log.Printf("@@@@@@@ ctrl retrieves priority: %v", priorities)
	if err != nil {
		return err
	}

	c.pub.Publish("update", tasks, priorities)
	return nil
}
