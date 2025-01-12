package router

import (
	"github.com/emrekentli/multitenant-boilerplate/src/handler"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(api fiber.Router) {
	productHandler := handler.NewProductHandler()
	productRoutes := api.Group("/products")
	productHandler.RegisterRoutes(productRoutes)
}
