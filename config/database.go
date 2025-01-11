package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type DatabaseConfig struct {
	DB          *pgxpool.Pool
	Host        string `yaml:"host" env:"DB_HOST"`
	Username    string `yaml:"username" env:"DB_USER"`
	Password    string `yaml:"password" env:"DB_PASS"`
	DBName      string `yaml:"db_name" env:"DB_NAME"`
	Port        int    `yaml:"port" env:"DB_PORT"`
	Connections int    `yaml:"connections" env:"DB_CONNECTIONS"` // Maksimum bağlantı sayısı
}

func (d *DatabaseConfig) InitDB() error {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		d.Username, d.Password, d.Host, d.Port, d.DBName,
	)

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
		return err
	}

	config.MaxConns = int32(d.Connections)    // Maksimum bağlantı sayısı
	config.MinConns = 5                       // Minimum bağlantı sayısı
	config.MaxConnLifetime = 30 * time.Minute // Bağlantıların maksimum yaşam süresi
	config.MaxConnIdleTime = 5 * time.Minute  // Boşta kalan bağlantılar için maksimum süre

	d.DB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to initialize database pool: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = d.DB.Ping(ctx); err != nil {
		log.Fatalf("Database connection ping failed: %v", err)
		return err
	}

	log.Println("Database connection established successfully.")
	return nil
}

func (d *DatabaseConfig) RunMigrations() error {
	if err := d.ApplyMigrations("migrations/public"); err != nil {
		return err
	}

	if err := d.ApplyMigrations("migrations/tenant"); err != nil {
		return err
	}

	log.Println("Migrations applied successfully.")
	return nil
}

func (d *DatabaseConfig) ApplyMigrations(migrationFolder string) error {
	files, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		return fmt.Errorf("failed to read migration folder %s: %v", migrationFolder, err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			if err := d.applySQLFile(migrationFolder + "/" + file.Name()); err != nil {
				return fmt.Errorf("failed to apply SQL file %s: %v", file.Name(), err)
			}
		}
	}

	return nil
}

func (d *DatabaseConfig) applySQLFile(filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %v", filePath, err)
	}

	sql := string(sqlBytes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = d.DB.Exec(ctx, sql)
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

func (d *DatabaseConfig) ExecuteQuery(ctx context.Context, query string, args ...interface{}) error {
	start := time.Now()
	_, err := d.DB.Exec(ctx, query, args...)
	duration := time.Since(start)

	logSQL(ctx, query, args, duration, err)
	return err
}

func (d *DatabaseConfig) Close() {
	if d.DB != nil {
		d.DB.Close()
		log.Println("Database connection pool closed.")
	}
}
