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
		c.eventLoop(ctx)
	}()
}

func (c *Ctrl) eventLoop(ctx context.Context) {
	for {
		// TODO: support cancel using context
		com := <-c.queue

		err := c.eventDispatch(ctx, com.event, com.data)
		if err != nil {
			c.errChan <- err
		}

		err = c.Update(context.Background())
		if err != nil {
			c.errChan <- err
		}
	}
}

func (c *Ctrl) eventDispatch(ctx context.Context, event delegateEvent, data []interface{}) error {
	var err error

	switch event {
	case AddTask:
		task := (*service.Task)(data[0].(*KeyedTask))
		task.IsDirty = true

		err = c.m.Create(ctx, task)
	case UpdateTask:
		task := (*service.Task)(data[0].(*KeyedTask))
		task.IsDirty = true

		err = c.m.Update(ctx, task)
	case DeleteTask:
		task := (*service.Task)(data[0].(*KeyedTask))
		task.IsDirty = true

		err = c.m.Delete(ctx, task.Id)
	case UpdatePriority:
		priority := data[0].(map[string]*service.Priority)
		for k := range priority {
			priority[k].IsDirty = true
		}

		_, err = c.m.UpdatePriority(ctx, priority)
	case UpdateBulk:
		task := (*service.Task)(data[0].(*KeyedTask))
		task.IsDirty = true
		priority := data[1].(map[string]*service.Priority)
		for k := range priority {
			priority[k].IsDirty = true
		}

		err = c.m.Update(ctx, task)
		if err != nil {
			c.errChan <- err
		}
		_, err = c.m.UpdatePriority(ctx, priority)
	default:
		panic(fmt.Sprintf("unknown delegateCommand: %v", event))
	}

	return err
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

	ktasks := make(map[string]*KeyedTask)
	for k, v := range tasks {
		ktasks[k] = (*KeyedTask)(v)
	}

	c.pub.Publish("update", ktasks, priorities)
	return nil
}
