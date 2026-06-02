package main

import (
	"context"
	"fmt"

	"gorm-unloaded/db"
)

func main() {
	db := db.MustNew()
	mustBootstrap(db)

	if err := db.Use(QueryPlugin{}); err != nil {
		panic(err)
	}

	// --- No tenant context — WHERE tenant_id filter is not injected ---
	var all []Product
	db.Find(&all)
	fmt.Printf("All products (no tenant): %d\n\n", len(all))

	// --- Tenant-A context — beforeQuery injects WHERE tenant_id = 'tenant-a' ---
	ctxA := context.WithValue(context.Background(), TenantIDKey, "tenant-a")
	var tenantAProducts []Product
	db.WithContext(ctxA).Find(&tenantAProducts)
	fmt.Printf("tenant-a products: %d\n\n", len(tenantAProducts))

	// --- Tenant-B context ---
	ctxB := context.WithValue(context.Background(), TenantIDKey, "tenant-b")
	var tenantBProducts []Product
	db.WithContext(ctxB).Find(&tenantBProducts)
	fmt.Printf("tenant-b products: %d\n\n", len(tenantBProducts))
}
