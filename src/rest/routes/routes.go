package routes

import (
	"github.com/emrekentli/multitenant-boilerplate/src/rest/middlewares"
	"github.com/gofiber/fiber/v2"
)

func LoadRoutes(app *fiber.App) {
	app.Use(middlewares.SetTenantContext)
	api := app.Group("api")
	web := app.Group("")
	ApiRoutes(api)
	WebRoutes(web)
}
