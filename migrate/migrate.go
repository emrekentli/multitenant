package migrate

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Migration işlemlerini başlat
func RunMigrations(db *sql.DB) error {
	// Migrations tablosunu kontrol et ve oluştur
	if err := ensureMigrationsTableExists(db, "public"); err != nil {
		return fmt.Errorf("public migrations table error: %v", err)
	}

	// Public migrationları uygula
	if err := applyMigrations(db, "public", "migrate/public"); err != nil {
		return fmt.Errorf("public migrations error: %v", err)
	}

	// Tenant'ları al
	tenants, err := getTenantsFromPublic(db)
	if err != nil {
		return fmt.Errorf("fetching tenants error: %v", err)
	}

	// Her tenant için migrationları uygula
	for _, tenant := range tenants {
		if err := applyMigrations(db, tenant, "migrate/tenant"); err != nil {
			return fmt.Errorf("tenant %s migration error: %v", tenant, err)
		}
	}

	log.Println("Migrations applied successfully.")
	return nil
}

// Migrationları uygula
func applyMigrations(db *sql.DB, tenant, folder string) error {
	// Şema ve migration tablosunu kontrol et ve oluştur
	if tenant != "public" {
		if err := ensureSchemaExists(db, tenant); err != nil {
			return fmt.Errorf("schema check error: %v", err)
		}
	}
	if err := ensureMigrationsTableExists(db, tenant); err != nil {
		return fmt.Errorf("migrations table error: %v", err)
	}

	// Migration dosyalarını uygula
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("reading folder %s error: %v", folder, err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			if err := applySQLFile(db, tenant, folder+"/"+file.Name()); err != nil {
				return fmt.Errorf("applying SQL file %s error: %v", file.Name(), err)
			}
		}
	}
	return nil
}

// Migration dosyasını uygula
func applySQLFile(db *sql.DB, tenant, filePath string) error {
	applied, err := isMigrationApplied(db, tenant, filePath)
	if err != nil {
		return err
	}
	if applied {
		log.Printf("Migration %s already applied for tenant %s, skipping.", filePath, tenant)
		return nil
	}

	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading SQL file %s error: %v", filePath, err)
	}

	replacedSQL := strings.ReplaceAll(string(sqlBytes), "schemaName", tenant)
	if _, err := db.Exec(replacedSQL); err != nil {
		return fmt.Errorf("executing SQL file %s error: %v", filePath, err)
	}

	if err := recordMigration(db, tenant, filePath); err != nil {
		return fmt.Errorf("recording migration %s error: %v", filePath, err)
	}

	log.Printf("Migration %s applied for tenant %s.", filePath, tenant)
	return nil
}

// Tenant'ları al
func getTenantsFromPublic(db *sql.DB) ([]string, error) {
	var tenants []string
	rows, err := db.Query("SELECT schema_name FROM public.tenants")
	if err != nil {
		return nil, fmt.Errorf("fetching tenants error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tenant string
		if err := rows.Scan(&tenant); err != nil {
			return nil, fmt.Errorf("scanning tenant error: %v", err)
		}
		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tenants error: %v", err)
	}
	return tenants, nil
}

// Şemayı kontrol et ve oluştur
func ensureSchemaExists(db *sql.DB, schema string) error {
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)`
	var exists bool
	if err := db.QueryRow(query, schema).Scan(&exists); err != nil {
		return fmt.Errorf("checking schema error: %v", err)
	}
	if !exists {
		if _, err := db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
			return fmt.Errorf("creating schema %s error: %v", schema, err)
		}
		log.Printf("Schema %s created successfully.", schema)
	}
	return nil
}

// Migration tablosunu kontrol et ve oluştur
func ensureMigrationsTableExists(db *sql.DB, schema string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`, schema)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("creating migrations table error: %v", err)
	}
	log.Printf("Migrations table ensured in schema %s.", schema)
	return nil
}

// Migration'ın daha önce uygulanıp uygulanmadığını kontrol et
func isMigrationApplied(db *sql.DB, tenant, filename string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.migrations WHERE filename = $1", tenant)
	var count int
	if err := db.QueryRow(query, filename).Scan(&count); err != nil {
		return false, fmt.Errorf("checking migration status error: %v", err)
	}
	return count > 0, nil
}

// Migration'ı kaydet
func recordMigration(db *sql.DB, tenant, filename string) error {
	query := fmt.Sprintf("INSERT INTO %s.migrations (filename) VALUES ($1)", tenant)
	if _, err := db.Exec(query, filename); err != nil {
		return fmt.Errorf("recording migration error: %v", err)
	}
	log.Printf("Recorded migration %s for tenant %s.", filename, tenant)
	return nil
}
