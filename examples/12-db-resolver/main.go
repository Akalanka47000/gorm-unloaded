package main

import (
	"fmt"
	"strings"
	"time"

	"gorm-unloaded/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	// MustNew creates the database on the primary and connects to it.
	primary := db.MustNew()
	mustBootstrap(primary)

	// Build the replica DSN pointing at the same database name.
	replicaDSN := strings.ReplaceAll(db.ReplicaDSN, "root", primary.Migrator().CurrentDatabase())

	// The replica runs with recovery_min_apply_delay=3s, so we wait for the
	// CREATE TABLE to replicate before registering the resolver.
	fmt.Println("Waiting for schema to replicate to replica...")
	time.Sleep(4 * time.Second)

	// Register the DB resolver plugin.
	// After this, reads (Find, First, Count) route to the replica automatically.
	// Writes (Create, Save, Update, Delete) continue to hit the primary.
	// db.Clauses(dbresolver.Write) and db.Clauses(dbresolver.Read) override routing per query.
	if err := primary.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
	})); err != nil {
		panic(fmt.Sprintf("failed to register dbresolver: %s", err))
	}

	// Write a new event to the primary.
	event := &Event{Title: "product launch"}
	err := primary.Create(event).Error
	if err != nil {
		panic(fmt.Sprintf("failed to create event: %s", err))
	}
	fmt.Printf("\nWrote %q to primary (id=%d)\n", event.Title, event.ID)

	// Poll the replica every second. The 3s replication delay means the event
	// won't appear immediately — this loop makes the lag visible.
	fmt.Println("\nPolling replica (reads auto-route there via dbresolver):")
	for range 6 {
		var count int64
		err := primary.Model(&Event{}).Count(&count).Error
		if err != nil {
			panic(fmt.Sprintf("failed to count events: %s", err))
		}
		fmt.Printf("  replica count: %d", count)
		if count > 0 {
			fmt.Println("  ← replica caught up")
			break
		}
		fmt.Println("  ← not yet replicated")
		time.Sleep(time.Second)
	}

	// Force a read against the primary to confirm the record is there.
	var primaryCount int64
	primary.Clauses(dbresolver.Write).Model(&Event{}).Count(&primaryCount)
	fmt.Printf("\nPrimary count (forced with dbresolver.Write): %d\n", primaryCount)
}
