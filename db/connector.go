package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MustNew is a convenience function for quickly setting up a temporary database connection
// without having to handle errors explicitly in the calling code to keep it shorter so more focus can be given to the main topic of the example.
//
// The function signature is pretty much the same as gorm.Open and works the same while doing the additional work of creating a temporary database for the caller's example.
// It accepts an optional gorm.Config parameter, allowing for flexible configuration of the database connection.
// If no configuration is provided, it defaults to an empty gorm.Config.
func MustNew(config ...*gorm.Config) *gorm.DB {
	dsn := strings.ReplaceAll(PrimaryDSN, "root", callerExampleName())
	err := createTemporaryDatabase(dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to create temporary database: %s", err))
	}
	if len(config) == 0 {
		config = append(config, &gorm.Config{})
	}
	db, err := gorm.Open(postgres.Open(dsn), config[0])
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err))
	}
	fmt.Printf("\n\033[32mConnected to database:\033[0m \033[1m%s\033[0m\n\n", db.Migrator().CurrentDatabase())
	return db
}

// createTemporaryDatabase creates a temporary database using the provided DSN.
// It parses the DSN to extract the database name and then connects to the PostgreSQL server to create a new database with that name.
func createTemporaryDatabase(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	u.Path = "/"

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Failed to close database connection: %s\n", err)
		}
	}()

	var exists bool
	if err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, dbName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		if _, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)); err != nil {
			return fmt.Errorf("failed to create test database: %w", err)
		}
	}

	return nil
}
