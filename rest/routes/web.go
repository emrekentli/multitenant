package routes

import (
	"github.com/gofiber/fiber/v2"
)

func WebRoutes(web fiber.Router) {
	web.Get("/", Home)
}

func Home(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello, World!")
}
