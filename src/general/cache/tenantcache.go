package cache

import (
	"app/src/general/database"
	"context"
	"sync"
)

var (
	domainToSchema map[string]string
	tenantMu       sync.RWMutex
)

func LoadTenantsToMemory() error {
	rows, err := database.DB.Query(context.Background(), "SELECT schema_name, domain FROM public.tenants")
	if err != nil {
		return err
	}
	defer rows.Close()

	tmp := make(map[string]string)
	for rows.Next() {
		var schema, domain string
		if err := rows.Scan(&schema, &domain); err != nil {
			return err
		}
		tmp[domain] = schema
	}

	tenantMu.Lock()
	domainToSchema = tmp
	tenantMu.Unlock()
	return nil
}

func GetSchemaByDomain(domain string) (string, bool) {
	tenantMu.RLock()
	defer tenantMu.RUnlock()
	schema, exists := domainToSchema[domain]
	return schema, exists
}
