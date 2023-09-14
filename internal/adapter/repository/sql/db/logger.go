package db

import (
	"context"
	"database/sql"
	"strings"
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

	subLogger := port.GetCtxSubLogger(ctx, h.logger)
	logEvent := subLogger.Info()
	if event.Err != nil {
		logEvent = subLogger.Error().Err(event.Err)
	}
	logEvent.
		Field("duration", duration).
		Field("query", strings.ReplaceAll(event.Query, "\"", "")).
		Msgf("SQL %s", event.Operation())
}
