package main

import (
	"context"

	hashirac "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
)

// Model represents model of hashira's mvc
type Model struct {
	hashirac *hashirac.Client
}

// SetHashiraClient sets hashira client
func (m *Model) SetHashiraClient(cli *hashirac.Client) {
	m.hashirac = cli
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
func (m *Model) UpdatePriority(ctx context.Context, priorities map[string]*service.Priority) (map[string]*service.Priority, error) {
	return m.hashirac.UpdatePriority(ctx, priorities)
}
