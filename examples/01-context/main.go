package main

import (
	"context"
	"fmt"
	"gorm-unloaded/db"
)

func main() {
	db := db.MustNew()

	ctx := context.WithValue(context.Background(), "key", "123")

	// This context value can be later accessed via db.Statement.Context
	// and serves multiple purposes, such as passing request-scoped values, deadlines, and cancellation signals across API boundaries and between processes.
	//
	// This can technically be inlined along with the main query, but for demonstration purposes, it is kept separate here to show that the context is indeed attached to the db instance and can be accessed later.
	db = db.WithContext(ctx)

	var result string

	err := db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		panic("failed to execute query")
	}

	fmt.Println("Query result:", result)

	fmt.Println("Context value retrieved from db statement -->", db.Statement.Context.Value("key")) // Outputs 123
}
