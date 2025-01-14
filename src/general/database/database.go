package database

import (
	"app/config"
	"app/src/general/util/secret_reader"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
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

func Connect() {
	if DB != nil {
		return
	}
	dbConfig, connString := getConfig()
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	ping(conn)
	DB = &Connection{conn, connString}
}

func Close() {
	log.Printf("Disconnected from database")
	DB.Close()
}

func getConfig() (*pgxpool.Config, string) {
	password := secret_reader.ReadSecret(config.DbPassword)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", config.DbUsername, password, config.DbHost, config.DbPort, config.DbDatabase, config.DbSslMode, config.DbSchema)
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Error parsing database config: %v", err)
	}
	return dbConfig, connectionString
}
func ping(conn *pgxpool.Pool) {
	err := conn.Ping(context.Background())
	if err != nil {
		log.Printf("Error pinging database: %v", err)
	} else {
		log.Printf("Connected to database")
	}
}
