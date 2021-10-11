package main

import (
	"context"
	"fmt"

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
}

func NewModel(hc *hashirac.Client, sc *syncutil.Client) *Model {
	return &Model{
		hashirac:   hc,
		syncclient: sc,
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
	if err := m.sync(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to sync: %w", err)
	}

	return p, nil
}

func (m *Model) Create(ctx context.Context, task *service.Task) error {
	if err := m.hashirac.Create(ctx, task); err != nil {
		return fmt.Errorf("failed to create a new task: %w", err)
	}
	if err := m.sync(context.Background()); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (m *Model) Update(ctx context.Context, task *service.Task) error {
	if err := m.hashirac.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update a task: %w", err)
	}
	if err := m.sync(context.Background()); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (m *Model) Delete(ctx context.Context, id string) error {
	if err := m.hashirac.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete a task: %w", err)
	}
	if err := m.sync(context.Background()); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (m *Model) SetAccessToken(accesstoken string) {
	m.accesstoken = accesstoken
}

func (m *Model) sync(ctx context.Context) error {
	if m.accesstoken == "" {
		return nil
	}
	if err := m.syncclient.Upload(m.accesstoken, syncutil.UploadDirtyOnly); err != nil {
		return fmt.Errorf("failed to upload tasks: %w", err)
	}
	if err := m.syncclient.Download(m.accesstoken); err != nil {
		return fmt.Errorf("failed to download tasks: %w", err)
	}
	return nil
}
