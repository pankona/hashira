package main

import (
	"context"
	"fmt"
	"log"

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
	event delegateEvent
	data  []interface{}
}

// Initialize initializes controller
func (c *Ctrl) Initialize() {
	c.queue = make(chan delegateCommand, 128)

	ctx := context.Background()

	go func() {
		// TODO: support cancel using context
		err := <-c.errChan
		log.Printf("[ERROR] %v", err)
	}()

	go func() {
		for {
			// TODO: support cancel using context
			com := <-c.queue

			event := com.event
			data := com.data

			var err error

			switch event {
			case AddTask:
				err = c.m.hashirac.Create(ctx, (*service.Task)(data[0].(*keyedTask)))
			case UpdateTask:
				err = c.m.hashirac.Update(ctx, (*service.Task)(data[0].(*keyedTask)))
			case DeleteTask:
				err = c.m.hashirac.Delete(ctx, (*service.Task)(data[0].(*keyedTask)).Id)
			case UpdatePriority:
				_, err = c.m.hashirac.UpdatePriority(ctx, data[0].([]*service.Priority))
			case UpdateBulk:
				err = c.m.hashirac.Update(ctx, (*service.Task)(data[0].(*keyedTask)))
				if err != nil {
					c.errChan <- err
				}
				_, err = c.m.hashirac.UpdatePriority(ctx, data[1].([]*service.Priority))
			default:
				panic(fmt.Sprintf("unknown delegateCommand: %v", event))
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
func (c *Ctrl) Delegate(event delegateEvent, data ...interface{}) (err error) {
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

	ktasks := make([]*keyedTask, len(tasks))
	for i, t := range tasks {
		ktasks[i] = (*keyedTask)(t)
	}

	c.pub.Publish("update", ktasks, priorities)
	return nil
}
