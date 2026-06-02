package main

import (
	"fmt"

	"gorm-unloaded/db"
)

// Status is stored as plain text in PostgreSQL.
// driver.Valuer controls what goes in; sql.Scanner controls what comes out.
// GORM calls both automatically — no extra configuration needed.
//
// For more complex column types (JSON, arrays, custom PostgreSQL types),
// the gorm.io/datatypes package provides ready-made implementations so you
// don't need to write Valuer/Scanner by hand:
//
//	datatypes.JSON        — stores a Go map/struct as a JSONB column
//	datatypes.Date        — date-only column without time component
//	datatypes.StringArray — PostgreSQL text[] column
//	datatypes.JSONSlice   — typed slice backed by JSONB

func main() {
	database := db.MustNew()
	mustBootstrap(database)

	// GORM calls Status.Value() when building the WHERE clause bind variable.
	var active []Order
	database.Where("status = ?", StatusActive).Find(&active)
	fmt.Printf("Active orders (%d):\n", len(active))
	for _, o := range active {
		fmt.Printf("  %s  status=%s\n", o.Reference, o.Status)
	}

	// Updating via the enum — Value() is called when GORM binds the new value.
	database.Model(&Order{}).
		Where("status = ?", StatusPending).
		Update("status", StatusActive)

	// Reading back — Scan() is called for every row returned.
	var all []Order
	database.Find(&all)
	fmt.Printf("\nAll orders after promoting pending → active:\n")
	for _, o := range all {
		fmt.Printf("  %s  status=%s\n", o.Reference, o.Status)
	}
}
