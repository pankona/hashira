package main

import (
	"context"

	"github.com/pankona/hashira/service"
)

// Ctrl represents controller of hashira-cui's mvc
type Ctrl struct {
	m       *Model
	pub     Publisher
	queue   chan delegateCommand
	errChan chan error
}

type delegateCommand struct {
	event string
	data  interface{}
}

// Initialize initializes controller
func (c *Ctrl) Initialize() {
	c.queue = make(chan delegateCommand, 128)

	go func() {
		for {
			// TODO: support cancel using context
			com := <-c.queue

			event := com.event
			data := com.data

			var err error
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

			if err != nil {
				c.errChan <- err
			}

			err = c.Update(context.Background())
			if err != nil {
				c.errChan <- err
			}
		}
	}()
}

// SetPublisher sets controller's Publisher
func (c *Ctrl) SetPublisher(p Publisher) {
	c.pub = p
}

// Delegate is called from view to delegate functionality that are not
// covered by view.
func (c *Ctrl) Delegate(event string, data interface{}) (err error) {
	c.queue <- delegateCommand{event: event, data: data}
	return nil
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
