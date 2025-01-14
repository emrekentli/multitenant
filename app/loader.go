package app

import (
	"app/config/server"
	"app/migrations"
	"app/src/general/database"
)

func Load() {
	database.Connect()
	defer database.Close()
	migrations.RunMigrations()
	server.New()
}
