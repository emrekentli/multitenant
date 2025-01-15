package migrations

import (
	"app/src/general/database"
	"app/src/general/util/hash"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func RunMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1130*time.Second)
	defer cancel()

	// Migrations tablosunu kontrol et ve oluştur
	if err := ensureMigrationsTableExists(ctx, "public"); err != nil {
		log.Printf("Error ensuring migrations table for public schema: %v", err)
		return fmt.Errorf("public migrations table error: %v", err)
	}

	// Public migrationları uygula
	if err := applyMigrations(ctx, "public", "migrations/public"); err != nil {
		log.Printf("Error applying migrations for public schema: %v", err)
		return fmt.Errorf("public migrations error: %v", err)
	}

	// Tenant'ları al
	tenants, err := getTenantsFromPublic(ctx)
	if err != nil {
		log.Printf("Error fetching tenants: %v", err)
		return fmt.Errorf("fetching tenants error: %v", err)
	}

	// Her tenant için migrationları uygula
	for _, tenant := range tenants {
		if err := applyMigrations(ctx, tenant, "migrations/tenant"); err != nil {
			log.Printf("Error applying migrations for tenant %s: %v", tenant, err)
			return fmt.Errorf("tenant %s migration error: %v", tenant, err)
		}
	}

	log.Println("Migrations applied successfully.")
	return nil
}

func applyMigrations(ctx context.Context, tenant, folder string) error {
	if tenant != "public" {
		if err := ensureSchemaExists(ctx, tenant); err != nil {
			log.Printf("Error ensuring schema for tenant %s: %v", tenant, err)
			return fmt.Errorf("schema check error: %v", err)
		}
	}
	if err := ensureMigrationsTableExists(ctx, tenant); err != nil {
		log.Printf("Error ensuring migrations table for tenant %s: %v", tenant, err)
		return fmt.Errorf("migrations table error: %v", err)
	}

	entries, err := os.ReadDir(folder)
	if err != nil {
		log.Printf("Error reading folder %s: %v", folder, err)
		return fmt.Errorf("reading folder %s error: %v", folder, err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sql") {
			if err := applySQLFile(ctx, tenant, folder+"/"+entry.Name()); err != nil {
				log.Printf("Error applying SQL file %s for tenant %s: %v", entry.Name(), tenant, err)
				return fmt.Errorf("applying SQL file %s error: %v", entry.Name(), err)
			}
		}
	}
	return nil
}

func applySQLFile(ctx context.Context, tenant, filePath string) error {
	applied, err := isMigrationApplied(ctx, tenant, filePath)
	if err != nil {
		log.Printf("Error checking if migration %s is applied for tenant %s: %v", filePath, tenant, err)
		return err
	}
	if applied {
		log.Printf("Migration %s already applied for tenant %s, skipping.", filePath, tenant)
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading SQL file %s: %v", filePath, err)
		return fmt.Errorf("reading SQL file %s error: %v", filePath, err)
	}
	defaultPassword := hash.Hash("123456")
	replacedSQL := strings.ReplaceAll(strings.ReplaceAll(string(data), "schemaName", tenant), "defaultPassword", defaultPassword)
	if _, err := database.DB.Exec(ctx, replacedSQL); err != nil {
		log.Printf("Error executing SQL file %s for tenant %s: %v", filePath, tenant, err)
		return fmt.Errorf("executing SQL file %s error: %v", filePath, err)
	}

	if err := recordMigration(ctx, tenant, filePath); err != nil {
		log.Printf("Error recording migration %s for tenant %s: %v", filePath, tenant, err)
		return fmt.Errorf("recording migration error: %v", err)
	}

	log.Printf("Migration %s applied for tenant %s.", filePath, tenant)
	return nil
}

func getTenantsFromPublic(ctx context.Context) ([]string, error) {
	var tenants []string
	rows, err := database.DB.Query(ctx, "SELECT schema_name FROM public.tenants")
	if err != nil {
		log.Printf("Error fetching tenants from public schema: %v", err)
		return nil, fmt.Errorf("fetching tenants error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tenant string
		if err := rows.Scan(&tenant); err != nil {
			log.Printf("Error scanning tenant: %v", err)
			return nil, fmt.Errorf("scanning tenant error: %v", err)
		}
		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over tenants rows: %v", err)
		return nil, fmt.Errorf("iterating tenants error: %v", err)
	}
	return tenants, nil
}

func ensureSchemaExists(ctx context.Context, schema string) error {
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)`
	var exists bool
	if err := database.DB.QueryRow(ctx, query, schema).Scan(&exists); err != nil {
		log.Printf("Error checking schema existence for %s: %v", schema, err)
		return fmt.Errorf("checking schema error: %v", err)
	}
	if !exists {
		if _, err := database.DB.Exec(ctx, fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
			log.Printf("Error creating schema %s: %v", schema, err)
			return fmt.Errorf("creating schema %s error: %v", schema, err)
		}
		log.Printf("Schema %s created successfully.", schema)
	}
	return nil
}

func ensureMigrationsTableExists(ctx context.Context, schema string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`, schema)
	if _, err := database.DB.Exec(ctx, query); err != nil {
		log.Printf("Error creating migrations table in schema %s: %v", schema, err)
		return fmt.Errorf("creating migrations table error: %v", err)
	}
	log.Printf("Migrations table ensured in schema %s.", schema)
	return nil
}

func isMigrationApplied(ctx context.Context, tenant, filename string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.migrations WHERE filename = $1", tenant)
	var count int
	if err := database.DB.QueryRow(ctx, query, filename).Scan(&count); err != nil {
		log.Printf("Error checking migration status for tenant %s: %v", tenant, err)
		return false, fmt.Errorf("checking migration status error: %v", err)
	}
	return count > 0, nil
}

func recordMigration(ctx context.Context, tenant, filename string) error {
	query := fmt.Sprintf("INSERT INTO %s.migrations (filename) VALUES ($1)", tenant)
	if _, err := database.DB.Exec(ctx, query, filename); err != nil {
		log.Printf("Error recording migration %s for tenant %s: %v", filename, tenant, err)
		return fmt.Errorf("recording migration error: %v", err)
	}
	log.Printf("Recorded migration %s for tenant %s.", filename, tenant)
	return nil
}
