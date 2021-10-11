package main

import (
	"context"

	hashirac "github.com/pankona/hashira/client"
	"github.com/pankona/hashira/service"
	"github.com/pankona/hashira/sync/syncutil"
)

// Model represents model of hashira's mvc
type Model struct {
	hashirac   *hashirac.Client
	syncclient *syncutil.Client
}

func NewModel(hc *hashirac.Client, sc *syncutil.Client) *Model {
	return &Model{
		hashirac:   hc,
		syncclient: sc,
	}
}

// SetHashiraClient sets hashira client
func (m *Model) SetHashiraClient(cli *hashirac.Client) {
	m.hashirac = cli
}

// List retrieves task list using hashira client
func (m *Model) List(
	ctx context.Context) (map[string]*service.Task, error) {
	return m.hashirac.Retrieve(ctx)
}

// RetrievePriority retrieves priorities using hashira client
func (m *Model) RetrievePriority(
	ctx context.Context) (map[string]*service.Priority, error) {
	return m.hashirac.RetrievePriority(ctx)
}

// UpdatePriority updates priorities using hashira client
func (m *Model) UpdatePriority(
	ctx context.Context,
	p map[string]*service.Priority) (map[string]*service.Priority, error) {
	return m.hashirac.UpdatePriority(ctx, p)
}

func (m *Model) Create(ctx context.Context, task *service.Task) error {
	return m.hashirac.Create(ctx, task)
}

func (m *Model) Update(ctx context.Context, task *service.Task) error {
	return m.hashirac.Update(ctx, task)
}

func (m *Model) Delete(ctx context.Context, id string) error {
	return m.hashirac.Delete(ctx, id)
}
