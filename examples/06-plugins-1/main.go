package main

import (
	"context"
	"errors"
	"fmt"

	"gorm-unloaded/db"
)

func main() {
	db := db.MustNew()
	mustBootstrap(db)

	// Register the plugin — all callbacks apply to every subsequent operation on this db instance.
	if err := db.Use(UserPlugin{}); err != nil {
		panic(err)
	}

	// --- BeforeCreate / AfterCreate ---
	// Same scenarios as example 05, but driven by the plugin rather than model methods.
	alice := &User{Name: "Alice", Email: "alice@example.com"}
	if err := db.Create(alice).Error; err != nil {
		panic(err)
	}
	fmt.Printf("\nAlice's role after create: %q\n", alice.Role)

	bob := &User{Name: "Bob", Email: "bob@example.com", Role: "admin"}
	if err := db.Create(bob).Error; err != nil {
		panic(err)
	}

	// --- BeforeDelete ---
	// db.AddError() inside the callback sets db.Error and aborts — same result as returning
	// an error from a hook, but the mechanism is different.
	err := db.Delete(bob).Error
	if errors.Is(err, ErrAdminProtected) {
		fmt.Printf("\nDelete aborted: %v\n", err)
	}

	if err := db.Delete(alice).Error; err != nil {
		panic(err)
	}

	// --- BeforeUpdate ---
	charlie := &User{Name: "Charlie", Email: "charlie@example.com"}
	if err := db.Create(charlie).Error; err != nil {
		panic(err)
	}

	memberCtx := context.WithValue(context.Background(), ActorRoleKey, "member")
	err = db.WithContext(memberCtx).Model(charlie).Updates(map[string]interface{}{"role": "admin"}).Error
	if errors.Is(err, ErrUnauthorizedRoleChange) {
		fmt.Printf("\nPromotion blocked: %v\n", err)
	}

	adminCtx := context.WithValue(context.Background(), ActorRoleKey, "admin")
	if err = db.WithContext(adminCtx).Model(charlie).Updates(map[string]interface{}{"role": "admin"}).Error; err != nil {
		panic(err)
	}
	db.First(charlie, charlie.ID)
	fmt.Printf("\nCharlie's role after admin-authorised update: %q\n", charlie.Role)
}
