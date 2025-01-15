package app

import (
	"app/config/server"
	"app/migrations"
	"app/src/general/database"
)

func Load() {
	database.Connect()
	defer database.Close()
	err := migrations.RunMigrations()
	if err != nil {
		return
	}
	server.New()
}
