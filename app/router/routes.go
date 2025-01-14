package router

import (
	"app/src/api"
	"github.com/gofiber/fiber/v3"
)

func RegisterFiberRoutes(app *fiber.App) {
	api.Register(app)
}
