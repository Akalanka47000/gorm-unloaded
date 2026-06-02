package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrAdminProtected = errors.New("admin users cannot be deleted")
var ErrUnauthorizedRoleChange = errors.New("only admins can grant the admin role")

type actorKey string

const ActorRoleKey actorKey = "actorRole"

// UserPlugin implements gorm.Plugin — the two-method interface GORM requires.
// Plugins register callbacks centrally via db.Use(), decoupling lifecycle behaviour
// from the model definition entirely. Hooks live on the model; plugins live outside it.
type UserPlugin struct{}

func (UserPlugin) Name() string { return "user_plugin" }

// Initialize is called once by db.Use(). Register all callbacks here.
// Callback keys are namespaced strings — "plugin_name:callback_name".
func (UserPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:create").Register("user_plugin:before_create", beforeCreate)
	db.Callback().Create().After("gorm:create").Register("user_plugin:after_create", afterCreate)
	db.Callback().Delete().Before("gorm:delete").Register("user_plugin:before_delete", beforeDelete)
	db.Callback().Update().Before("gorm:update").Register("user_plugin:before_update", beforeUpdate)
	return nil
}

// Plugin callbacks receive *gorm.DB — but unlike hooks, the signature has no error return.
// Errors are communicated via db.AddError(err), which sets db.Error and aborts the chain.

func beforeCreate(db *gorm.DB) {
	user, ok := db.Statement.Model.(*User)
	if !ok || user.Role != "" {
		return
	}
	user.Role = "member"
	fmt.Printf("\n[plugin:BeforeCreate] %q → role defaulted to %q\n", user.Name, user.Role)
}

func afterCreate(db *gorm.DB) {
	user, ok := db.Statement.Model.(*User)
	if !ok {
		return
	}
	// db.Statement.Schema is the fully parsed schema — Table, Fields, Relationships, Indices.
	// It is populated during gorm:initialize and available in all subsequent callbacks.
	fmt.Printf("\n[plugin:AfterCreate]  %q created (id=%d, table=%q, schema_fields=%d)\n",
		user.Name, user.ID, db.Statement.Schema.Table, len(db.Statement.Schema.Fields))
}

func beforeDelete(db *gorm.DB) {
	user, ok := db.Statement.Model.(*User)
	if !ok || user.Role != "admin" {
		return
	}
	fmt.Printf("\n[plugin:BeforeDelete] blocked — %q is an admin\n", user.Name)
	db.AddError(ErrAdminProtected)
}

func beforeUpdate(db *gorm.DB) {
	if !db.Statement.Changed("Role") {
		return
	}
	dest, ok := db.Statement.Dest.(map[string]interface{})
	if !ok {
		return
	}
	newRole, _ := dest["role"].(string)
	actorRole, _ := db.Statement.Context.Value(ActorRoleKey).(string)

	fmt.Printf("\n[plugin:BeforeUpdate] table=%q  Role→%q  actor=%q\n",
		db.Statement.Table, newRole, actorRole)

	if newRole == "admin" && actorRole != "admin" {
		db.AddError(ErrUnauthorizedRoleChange)
	}
}
