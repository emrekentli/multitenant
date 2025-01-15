package router

import (
	"app/src/api"
	"app/src/general/util/rest"
	"github.com/gofiber/fiber/v3"
)

func RegisterFiberRoutes(app *fiber.App) {
	api.Register(app)
	app.Use(func(c fiber.Ctx) error {
		return rest.ErrorRes(c, "404", "No route found")
	})
}
