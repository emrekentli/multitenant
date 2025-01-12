package app

import (
	"github.com/emrekentli/multitenant-boilerplate/config"
	"github.com/emrekentli/multitenant-boilerplate/migrations"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

var Http *config.AppConfig

func Load() {
	Http = &config.AppConfig{}
	Http.Setup()
	LoadBuiltInMiddlewares(Http)
	err := migrations.RunMigrations()
	if err != nil {
		return
	}
}

func LoadBuiltInMiddlewares(app *config.AppConfig) {
	app.Server.Use(recover.New())
	app.Server.Use(etag.New())
	app.Server.Use(compress.New(compress.Config{
		Level: 1, // Sıkıştırma seviyesi (1: minimum, 9: maksimum)
	}))

}
