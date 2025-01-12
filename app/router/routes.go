package router

import (
	"github.com/emrekentli/multitenant-boilerplate/src/api"
	"github.com/gofiber/fiber/v3"
)

func LoadRoutes(app *fiber.App) {
	api.Register(app)
}
