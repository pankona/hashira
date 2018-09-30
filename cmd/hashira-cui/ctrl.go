package main

import (
	"context"
	"fmt"

	"github.com/pankona/hashira/service"
	"github.com/pkg/errors"
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

func (c *Ctrl) Delegate(event string, data interface{}) (err error) {
	defer func() {
		e := c.Update(context.Background())
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("failed to Update: %s", e.Error()))
		}
	}()

	ctx := context.Background()

	switch event {
	case "add":
		err = c.m.hashirac.Create(ctx, data.(*service.Task))
	case "update":
		err = c.m.hashirac.Update(ctx, data.(*service.Task))
	case "delete":
		err = c.m.hashirac.Delete(ctx, data.(*service.Task).Id)
	case "updatePriority":
		_, err = c.m.hashirac.UpdatePriority(ctx, data.([]*service.Priority))
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
