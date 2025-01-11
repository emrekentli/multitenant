package config

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type DatabaseDriver struct {
	Driver      string `yaml:"driver" env:"DB_DRIVER"`
	Host        string `yaml:"host" env:"DB_HOST"`
	Username    string `yaml:"username" env:"DB_USER"`
	Password    string `yaml:"password" env:"DB_PASS"`
	DBName      string `yaml:"db_name" env:"DB_NAME"`
	Port        int    `yaml:"port" env:"DB_PORT"`
	Connections int    `yaml:"connections" env:"DB_CONNECTIONS"`
}

type DatabaseConfig struct {
	DB      *sql.DB
	Drivers map[string]DatabaseDriver `yaml:"drivers"`
	Default DatabaseDriver            `yaml:"default" env:"DEFAULT_DB_DRIVER"`
}

func (d *DatabaseConfig) Setup() error {
	var err error
	var connectionString string

	if d.DB != nil {
		return nil
	}

	switch d.Default.Driver {
	case "postgres":
		connectionString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			d.Default.Host, d.Default.Port, d.Default.Username, d.Default.DBName, d.Default.Password)
		d.DB, err = sql.Open("postgres", connectionString)

	default:
		return fmt.Errorf("unsupported database driver: %s", d.Default.Driver)
	}

	if err != nil {
		log.Fatal("Database connection failed:", err)
		return err
	}

	d.DB.SetMaxOpenConns(d.Default.Connections)
	d.DB.SetMaxIdleConns(d.Default.Connections)
	d.DB.SetConnMaxLifetime(24 * time.Hour)

	return nil
}
func (d *DatabaseConfig) RunMigrations() error {
	// public ve tenant klasörlerindeki SQL dosyalarını oku
	err := d.ApplyMigrations("migrations/public")
	if err != nil {
		return err
	}

	err = d.ApplyMigrations("migrations/tenant")
	if err != nil {
		return err
	}

	log.Println("Migrations applied successfully.")
	return nil
}

// SQL dosyalarını uygula
func (d *DatabaseConfig) ApplyMigrations(migrationFolder string) error {
	files, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		return fmt.Errorf("failed to read migration folder %s: %v", migrationFolder, err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			err := d.applySQLFile(migrationFolder + "/" + file.Name())
			if err != nil {
				return fmt.Errorf("failed to apply SQL file %s: %v", file.Name(), err)
			}
		}
	}

	return nil
}

// SQL dosyasını çalıştır
func (d *DatabaseConfig) applySQLFile(filePath string) error {
	// SQL dosyasını oku
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %v", filePath, err)
	}

	// SQL sorgusunu çalıştır
	sql := string(sqlBytes)
	_, err = d.DB.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file %s: %v", filePath, err)
	}

	log.Printf("Successfully applied migration from file: %s", filePath)
	return nil
}

func logSQL(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	if err != nil {
		log.Printf("ERROR: %v | Query: %s | Args: %v | Duration: %v\n", err, query, args, duration)
		return
	}
	log.Printf("INFO: Query: %s | Args: %v | Duration: %v\n", query, args, duration)
}

// Execute a query example
func (d *DatabaseConfig) ExecuteQuery(ctx context.Context, query string, args ...interface{}) error {
	start := time.Now()
	_, err := d.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	logSQL(ctx, query, args, duration, err)
	return err
}
