package dissector

import (
	"cocoon/pkg/db"
	"context"
	"go.uber.org/zap"
)

type DissectManager struct {
	logger *zap.Logger
	ctx    context.Context
}

func NewDissectManager(ctx context.Context, logger *zap.Logger) *DissectManager {
	return &DissectManager{
		logger: logger,
		ctx:    ctx,
	}
}

func (d *DissectManager) Dissect(session string, database *db.Database) error {
	cursor, err := database.ReadTcpTrafficBySession(session)
	if err != nil {
		return err
	}
	worker := NewDissectWorker(d.ctx, d.logger, cursor)
	return worker.Start()
}
