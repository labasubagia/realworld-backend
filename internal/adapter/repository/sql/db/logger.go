package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/port"
	"github.com/uptrace/bun"
)

type LoggerHook struct {
	verbose bool
	logger  port.Logger
}

func (h *LoggerHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *LoggerHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	if !h.verbose {
		switch event.Err {
		case nil, sql.ErrNoRows, sql.ErrTxDone:
			return
		}
	}
	now := time.Now()
	duration := now.Sub(event.StartTime)

	logger := h.logger.Info()
	if event.Err != nil {
		logger = h.logger.Error().Err(event.Err)
	}
	logger.
		Field("duration", duration).
		Field("query", event.Query).
		Msgf("SQL %s", event.Operation())
}
