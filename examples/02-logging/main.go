package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm-unloaded/db"
	"gorm.io/gorm"
)

func main() {
	db := db.MustNew()

	executeQueries := func(ctx context.Context, db *gorm.DB) {
		var result string

		err := db.WithContext(ctx).Raw("SELECT 1 FROM (SELECT pg_sleep(1)) AS t").Scan(&result).Error
		if err != nil {
			panic("failed to execute query")
		}
		_ = db.WithContext(ctx).Raw("SELECT pg_unknown_command(1)").Scan(&result).Error // This will cause an error and trigger an error log
	}

	// --- With the default logger....

	executeQueries(context.Background(), db)

	// --- With custom logger enriched with context....

	db = db.Session(&gorm.Session{
		Logger: GormLogger(), // Using the custom logger defined in logger.go
	})

	ctx := context.Background()
	ctx = log.Logger.With().Str("requestID", "abc123").Logger().WithContext(ctx) // Enriching the context with a request ID for better traceability in logs

	executeQueries(ctx, db)

	// Logger can also be set while opening the database connection using the optional gorm.Config parameter in db.MustNew, but for demonstration purposes, it is set here to show that it can be changed at any point in the code.
	// db = db.MustNew(&gorm.Config{
	// 	Logger: GormLogger(),
	// })

	// --- Force logging a single query via .Debug() method ---

	_ = db.WithContext(ctx).Debug().Exec("SELECT 2")
}
