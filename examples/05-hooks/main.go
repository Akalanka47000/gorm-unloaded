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

	// --- BeforeCreate: default field injection ---
	// Role is intentionally left empty — BeforeCreate fills it in before the INSERT.
	alice := &User{Name: "Alice", Email: "alice@example.com"}
	if err := db.Create(alice).Error; err != nil {
		panic(err)
	}
	fmt.Printf("\nAlice's role after create: %q\n\n", alice.Role)

	// --- AfterCreate: ID is populated before the hook runs ---
	// Bob is created with an explicit role. AfterCreate shows the assigned ID.
	bob := &User{Name: "Bob", Email: "bob@example.com", Role: "admin"}
	if err := db.Create(bob).Error; err != nil {
		panic(err)
	}

	// --- BeforeDelete: returning an error aborts the operation ---
	// Attempting to delete an admin user is blocked by the hook.
	err := db.Delete(bob).Error
	if errors.Is(err, ErrAdminProtected) {
		fmt.Printf("Delete was aborted: %v\n", err)
	}

	// Regular member deletion proceeds — BeforeDelete returns nil.
	if err := db.Delete(alice).Error; err != nil {
		panic(err)
	}

	// --- BeforeUpdate: tx.Statement field inspection ---
	// BeforeUpdate uses tx.Statement.Changed() to detect whether Role is part of
	// this specific update, tx.Statement.Dest to read the incoming value, and
	// tx.Statement.Context to check who is making the change.
	charlie := &User{Name: "Charlie", Email: "charlie@example.com"}
	if err := db.Create(charlie).Error; err != nil {
		panic(err)
	}

	// A member-level actor tries to grant admin — BeforeUpdate blocks it.
	memberCtx := context.WithValue(context.Background(), ActorRoleKey, "member")
	err = db.WithContext(memberCtx).Model(charlie).Updates(map[string]interface{}{
		"role": "admin",
	}).Error
	if errors.Is(err, ErrUnauthorizedRoleChange) {
		fmt.Printf("\nPromotion blocked: %v\n", err)
	}

	// An admin actor makes the same change — BeforeUpdate allows it.
	adminCtx := context.WithValue(context.Background(), ActorRoleKey, "admin")
	if err = db.WithContext(adminCtx).Model(charlie).Updates(map[string]interface{}{
		"role": "admin",
	}).Error; err != nil {
		panic(err)
	}
	err = db.First(charlie, charlie.ID).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCharlie's role after admin-authorised update: %q\n", charlie.Role)
}
