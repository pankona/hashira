package main

import (
	"context"
	"fmt"

	"github.com/pankona/hashira/service"
	"github.com/pkg/errors"
)

// Ctrl represents controller of hashira-cui's mvc
type Ctrl struct {
	m   *Model
	pub Publisher
}

// Initialize initializes controller
func (c *Ctrl) Initialize() {}

// SetPublisher sets controller's Publisher
func (c *Ctrl) SetPublisher(p Publisher) {
	c.pub = p
}

// Delegate is called from view to delegate functionality that are not
// covered by view.
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

// Update retrieve latest information from model and
// reflect them to view via PubSub.
func (c *Ctrl) Update(ctx context.Context) error {
	// TODO:
	// List and RetrievePriority should archive in
	// one communication for better performance
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
