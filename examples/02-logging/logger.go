package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	SourceField    = "file"
	QueryField     = "query"
	DurationField  = "duration"
	SlowQueryField = "slowQuery"
	RowsField      = "rows"
)

// GormLogger creates a new logger for gorm.io/gorm
func GormLogger() gormLogger { // nolint:revive
	l := gormLogger{
		ignoreRecordNotFoundError: false,
		slowThreshold:             time.Millisecond * 300,
	}
	return l
}

type gormLogger struct {
	ignoreTrace               bool
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
	logLevel                  gormlogger.LogLevel
}

// LogMode log mode
func (l gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.logLevel = level
	return l
}

// Info logs info
func (l gormLogger) Info(ctx context.Context, format string, args ...any) {
	log.Ctx(ctx).Info().Msgf(format, args...)
}

// Warn logs warn messages
func (l gormLogger) Warn(ctx context.Context, format string, args ...any) {
	log.Ctx(ctx).Warn().Msgf(format, args...)
}

// Error logs error messages
func (l gormLogger) Error(ctx context.Context, format string, args ...any) {
	log.Ctx(ctx).Error().Msgf(format, args...)
}

// Trace logs sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.ignoreTrace {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		sql, rows := fc()
		log.Ctx(ctx).Error().Err(err).
			Str(QueryField, sql).Dur(DurationField, elapsed).
			Int64(RowsField, rows).Str(SourceField, utils.FileWithLineNum()).
			Msg("SQL query error")
	case l.slowThreshold != 0 && elapsed > l.slowThreshold:
		sql, rows := fc()
		log.Ctx(ctx).Warn().Bool(SlowQueryField, true).
			Str(QueryField, sql).Dur(DurationField, elapsed).
			Int64(RowsField, rows).Str(SourceField, utils.FileWithLineNum()).
			Msg("Slow SQL query")
	case l.logLevel == gormlogger.Info:
		sql, rows := fc()
		log.Ctx(ctx).Info().Str(QueryField, sql).
			Dur(DurationField, elapsed).
			Int64(RowsField, rows).Str(SourceField, utils.FileWithLineNum()).
			Msg("SQL query executed")
	}
}
