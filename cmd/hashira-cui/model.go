package main

import (
	"context"

	hashirac "github.com/pankona/hashira/hashira/client"
	"github.com/pankona/hashira/service"
)

type Model struct {
	hashirac *hashirac.Client
}

func (m *Model) SetHashiraClient(cli *hashirac.Client) {
	m.hashirac = cli
}

func (m *Model) List(ctx context.Context) ([]*service.Task, error) {
	return m.hashirac.Retrieve(ctx)
}
