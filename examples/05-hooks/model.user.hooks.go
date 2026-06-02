package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ErrAdminProtected is returned by BeforeDelete to abort deletion of admin users.
var ErrAdminProtected = errors.New("admin users cannot be deleted")

// ErrUnauthorizedRoleChange is returned by BeforeUpdate when a non-admin tries to grant admin.
var ErrUnauthorizedRoleChange = errors.New("only admins can grant the admin role")

type actorKey string

// ActorRoleKey is the context key for the role of the user performing the operation.
const ActorRoleKey actorKey = "actorRole"

// BeforeCreate fires before INSERT. Mutations to the receiver are persisted.
// Returning a non-nil error aborts the operation — no row is written.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "member" // inject default before the INSERT runs
	}
	fmt.Printf("[BeforeCreate] %q → role set to %q\n", u.Name, u.Role)
	return nil
}

// AfterCreate fires after a successful INSERT. The primary key is already populated
// on the receiver at this point, so u.ID is safe to read.
func (u *User) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("\n[AfterCreate]  %q created with id=%d\n", u.Name, u.ID)
	return nil
}

// BeforeDelete fires before DELETE. Returning an error cancels the deletion —
// no SQL is executed and the error is surfaced on db.Error.
func (u *User) BeforeDelete(tx *gorm.DB) error {
	if u.Role == "admin" {
		fmt.Printf("\n[BeforeDelete] blocked — %q is an admin\n\n", u.Name)
		return ErrAdminProtected
	}
	fmt.Printf("\n[BeforeDelete] %q approved for deletion\n\n", u.Name)
	return nil
}

// BeforeUpdate demonstrates reading multiple tx.Statement fields to enforce a
// context-aware authorization rule.
//
//	tx.Statement.Changed("Field") — true only if that field is part of this update
//	tx.Statement.Dest            — the map or struct carrying the incoming values
//	tx.Statement.Context         — the context attached via db.WithContext()
//	tx.Statement.Table           — the resolved table name at hook execution time
//
// Together they let the hook make a decision that depends on what is changing,
// what it is changing to, who is making the change, and which table is involved —
// all without the caller passing any extra arguments.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if !tx.Statement.Changed("Role") {
		return nil // Role is not part of this update — nothing to enforce
	}

	// tx.Statement.Dest holds what is being written.
	// When using db.Model().Update("Role", value), GORM stores the update as a map.
	dest, ok := tx.Statement.Dest.(map[string]interface{})
	if !ok {
		return nil
	}
	newRole, _ := dest["role"].(string)

	// tx.Statement.Context carries the actor identity set upstream via db.WithContext(ctx).
	actorRole, _ := tx.Statement.Context.Value(ActorRoleKey).(string)

	fmt.Printf("\n[BeforeUpdate] table=%q  Role → %q  (actor_role: %q)\n",
		tx.Statement.Table, newRole, actorRole)

	if newRole == "admin" && actorRole != "admin" {
		return ErrUnauthorizedRoleChange
	}
	return nil
}
