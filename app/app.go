package app

import (
	"github.com/emrekentli/multitenant-boilerplate/config"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var Http *config.AppConfig

func Load(configFile string) {
	Http = &config.AppConfig{ConfigFile: configFile}
	Http.Setup()
	LoadBuiltInMiddlewares(Http)
}

func LoadBuiltInMiddlewares(app *config.AppConfig) {
	app.Server.Use(recover.New())
	app.Server.Use(etag.New())
	app.Server.Use(compress.New(compress.Config{
		Level: 1,
	}))
	if app.Server.Debug {
		app.Server.Use(pprof.New())
	}
}
