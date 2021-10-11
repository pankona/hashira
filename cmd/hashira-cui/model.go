package main

import (
	"context"
	"fmt"
	"log"
	"time"

	hashirac "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
	"github.com/pankona/hashira/sync/syncutil"
)

// Model represents model of hashira's mvc
type Model struct {
	hashirac   *hashirac.Client
	syncclient *syncutil.Client

	// TODO: remove if accesstoken is held by syncclient
	accesstoken string
	syncChan    chan struct{}
}

func NewModel(hc *hashirac.Client, sc *syncutil.Client) *Model {
	return &Model{
		hashirac:   hc,
		syncclient: sc,
		syncChan:   make(chan struct{}),
	}
}

// List retrieves task list using hashira client
func (m *Model) List(ctx context.Context) (map[string]*service.Task, error) {
	return m.hashirac.Retrieve(ctx)
}

// RetrievePriority retrieves priorities using hashira client
func (m *Model) RetrievePriority(ctx context.Context) (map[string]*service.Priority, error) {
	return m.hashirac.RetrievePriority(ctx)
}

// UpdatePriority updates priorities using hashira client
func (m *Model) UpdatePriority(ctx context.Context, p map[string]*service.Priority) (map[string]*service.Priority, error) {
	p, err := m.hashirac.UpdatePriority(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to update priority: %w", err)
	}
	m.NotifySync()

	return p, nil
}

func (m *Model) Create(ctx context.Context, task *service.Task) error {
	if err := m.hashirac.Create(ctx, task); err != nil {
		return fmt.Errorf("failed to create a new task: %w", err)
	}
	m.NotifySync()

	return nil
}

func (m *Model) Update(ctx context.Context, task *service.Task) error {
	if err := m.hashirac.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update a task: %w", err)
	}
	m.NotifySync()

	return nil
}

func (m *Model) Delete(ctx context.Context, id string) error {
	if err := m.hashirac.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete a task: %w", err)
	}
	m.NotifySync()

	return nil
}

func (m *Model) SetAccessToken(accesstoken string) {
	m.accesstoken = accesstoken
}

func (m *Model) NotifySync() {
	select {
	case m.syncChan <- struct{}{}:
	default:
	}
}

func (m *Model) SyncOnNotify(ctx context.Context) error {
	if m.accesstoken == "" {
		return nil
	}

	var cancelFunc context.CancelFunc

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-m.syncChan:
			if cancelFunc != nil {
				cancelFunc()
			}

			ctx, cancel := context.WithCancel(ctx)
			cancelFunc = cancel

			go func(ctx context.Context) {
				select {
				case <-ctx.Done():
					// do nothing
				case <-time.After(1 * time.Minute):
					if err := m.syncclient.Upload(m.accesstoken, syncutil.UploadDirtyOnly); err != nil {
						log.Printf("failed to upload tasks: %v", err)
					}
					if err := m.syncclient.Download(m.accesstoken); err != nil {
						log.Printf("failed to download tasks: %v", err)
					}

				}
			}(ctx)
		}
	}
}
