package app

import (
	"app/config/server"
	"app/migrations"
	"app/src/general/cache"
	"app/src/general/database"
	"github.com/gofiber/fiber/v3/log"
)

func Load() {
	database.Connect()
	defer database.Close()
	err := migrations.RunMigrations()
	if err != nil {
		return
	}
	if err := cache.LoadTenantsToMemory(); err != nil {
		log.Fatalf("Tenant listesi yüklenemedi: %v", err)
	}
	server.New()
}
