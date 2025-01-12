package database

import (
	"context"
	"fmt"
	"github.com/emrekentli/multitenant-boilerplate/src/util/function"
	"github.com/emrekentli/multitenant-boilerplate/src/util/secret_reader"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Connection struct {
	*pgxpool.Pool
	ConnectionString string
}
type PgError struct {
	Code string
}

var (
	DB *Connection
)

type DatabaseConfig struct {
	Host        string
	Username    string
	Password    string
	DBName      string
	Port        string
	Connections string
	SSLMode     string
}

func (d *DatabaseConfig) InitDB() error {
	dbConfig, connString := d.getConfig()
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		log.Fatalf("Failed to initialize database pool: %v", err)
		return err
	}
	ping(conn)

	DB = &Connection{conn, connString}
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

	_, err = DB.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file %s: %v", filePath, err)
	}

	log.Printf("Successfully applied migration from file: %s", filePath)
	return nil
}

func (d *DatabaseConfig) Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection pool closed.")
	}
}

func (d *DatabaseConfig) getConfig() (*pgxpool.Config, string) {
	d.Host = os.Getenv("DB_HOST")
	d.Username = os.Getenv("DB_USERNAME")
	d.Password = os.Getenv("DB_PASSWORD")
	d.DBName = os.Getenv("DB_NAME")
	d.SSLMode = os.Getenv("DB_SSL_MODE")
	d.Port = os.Getenv("DB_PORT")
	d.Connections = os.Getenv("DB_CONNECTIONS")
	password := secret_reader.ReadSecret(d.Password)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", d.Username, password, d.Host, d.Port, d.DBName, d.SSLMode)
	config, err := pgxpool.ParseConfig(connectionString)
	config.MaxConns = function.StringToInt(d.Connections)
	config.MinConns = 5                       // Minimum number of connections
	config.MaxConnLifetime = 30 * time.Minute // Maximum connection lifetime
	config.MaxConnIdleTime = 5 * time.Minute  // Maximum idle connection time
	if err != nil {
		log.Printf("Error parsing database config: %v", err)
	}
	return config, connectionString
}
func ping(conn *pgxpool.Pool) {
	err := conn.Ping(context.Background())
	if err != nil {
		log.Printf("Error pinging database: %v", err)
	} else {
		log.Printf("Connected to database")
	}
}
