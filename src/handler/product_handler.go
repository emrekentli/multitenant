package handler

import (
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/emrekentli/multitenant-boilerplate/src/controller"
	"github.com/emrekentli/multitenant-boilerplate/src/repository"
	"github.com/emrekentli/multitenant-boilerplate/src/service"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductController *controller.ProductController
}

func NewProductHandler() *ProductHandler {
	productRepository := repository.NewProductRepository(app.Http.Database.DB)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)
	return &ProductHandler{ProductController: productController}
}

func (h *ProductHandler) RegisterRoutes(router fiber.Router) {
	router.Group("/products")
	router.Get("/", h.ProductController.GetAllProducts)
	router.Post("/", h.ProductController.CreateProduct)
	router.Get("/:id", h.ProductController.GetProductByID)
	router.Put("/:id", h.ProductController.UpdateProduct)
	router.Delete("/:id", h.ProductController.DeleteProduct)
}
