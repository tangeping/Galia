package mongodb

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type Model struct {
	cli *qmgo.QmgoClient
}

// NewModel returns a Model.
func NewModel(cfg *qmgo.Config, opt ...options.ClientOptions) (*Model, error) {
	client, err := qmgo.Open(context.Background(), cfg, opt...)
	if err != nil {
		return nil, err
	}
	return &Model{
		cli: client,
	}, nil
}

func (m *Model) Version() string {
	if m.cli == nil {
		return ""
	}
	return m.cli.ServerVersion()
}
