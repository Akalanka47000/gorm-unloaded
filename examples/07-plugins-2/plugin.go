package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tenantKey string

const TenantIDKey tenantKey = "tenantID"

// QueryPlugin demonstrates all injection points available in GORM's query callback chain.
//
// The chain has three named callbacks executed in order:
//
//	"gorm:query"       — builds the SQL (BuildQuerySQL) then executes it via ConnPool.QueryContext
//	"gorm:preload"     — loads associations registered with db.Preload(...)
//	"gorm:after_query" — calls AfterFind hooks on the loaded results
//
// Available injection points:
//
//	Before("gorm:query")       — fires before SQL is built; full access to modify the Statement
//	                             (add WHERE/ORDER/SELECT, read context, set up timing, abort early)
//	After("gorm:query")        — fires after rows are scanned into Dest; results are ready
//	                             (read SQL string, row count, duration, transform results)
//	Before("gorm:preload")     — same timing as After("gorm:query"), before associations load
//	After("gorm:preload")      — associations are now loaded; can inspect or modify them
//	Before("gorm:after_query") — same timing as After("gorm:preload"), before AfterFind fires
//	After("gorm:after_query")  — everything is done, AfterFind hooks have run
//	Replace("gorm:query", fn)  — swap out the entire query executor with your own implementation
//	Remove("gorm:query")       — remove a named callback from the chain entirely
//
// State between Before and After: use db.InstanceSet / db.InstanceGet.
// These write into db.Statement.Settings (a sync.Map) scoped to the current operation.
type QueryPlugin struct{}

func (QueryPlugin) Name() string { return "query_plugin" }

func (QueryPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("query_plugin:before", beforeQuery)
	db.Callback().Query().After("gorm:query").Register("query_plugin:after", afterQuery)
	return nil
}

// beforeQuery fires before BuildQuerySQL runs — the Statement has no SQL yet.
// This is the right place to modify WHERE, SELECT, ORDER, or any other clause.
func beforeQuery(db *gorm.DB) {
	// Pass timing state to afterQuery via InstanceSet — scoped to this operation only.
	db.InstanceSet("query_plugin:start", time.Now())

	// Read tenant from context and inject into WHERE via AddClause.
	// This happens transparently — callers only set the context, not the filter.
	tenantID, ok := db.Statement.Context.Value(TenantIDKey).(string)
	if !ok || tenantID == "" {
		return
	}
	db.Statement.AddClause(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{Column: "tenant_id", Value: tenantID},
		},
	})
	fmt.Printf("[plugin:before] injected WHERE tenant_id = %q\n", tenantID)
}

// afterQuery fires after rows are scanned into db.Statement.Dest.
// db.Statement.SQL.String() holds the exact SQL that was sent to the database.
// db.RowsAffected holds the number of rows returned.
func afterQuery(db *gorm.DB) {
	start, ok := db.InstanceGet("query_plugin:start")
	if !ok {
		return
	}
	elapsed := time.Since(start.(time.Time))
	fmt.Printf("[plugin:after]  sql=%q  rows=%d  duration=%s\n\n",
		db.Statement.SQL.String(), db.RowsAffected, elapsed.Round(time.Millisecond))
}
