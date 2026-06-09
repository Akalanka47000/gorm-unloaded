package main

import (
	"fmt"

	"gorm-unloaded/db"
)

func main() {
	database := db.MustNew()
	mustBootstrap(database)

	// The serializer means GORM and the DB see two different types for the same field.
	// Read via GORM — Go side receives a key=value string.
	var users []User
	if err := database.Find(&users).Error; err != nil {
		panic(err)
	}

	// Read via raw SQL — DB side shows the actual JSONB map stored on disk.
	type rawRow struct{ Name, Preferences string }
	var raw []rawRow
	err := database.Raw("SELECT name, preferences FROM users").Scan(&raw).Error
	if err != nil {
		panic("failed to fetch raw users")
	}

	fmt.Printf("  %-8s  %-40s  %s\n", "name", "Go (string)", "DB (JSONB)")
	fmt.Printf("  %-8s  %-40s  %s\n", "----", "----------", "---------")
	for i, u := range users {
		fmt.Printf("  %-8s  %-40s  %s\n", u.Name, u.Preferences, raw[i].Preferences)
	}
}
