package controller

import (
	"github.com/emrekentli/multitenant-boilerplate/app/middlewares"
	"github.com/emrekentli/multitenant-boilerplate/src/service"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{ProductService: productService}
}

func (h *ProductController) GetAllProducts(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	if tenantContext == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant context is missing"})
	}

	products, err := h.ProductService.GetAllProducts(tenantContext.SchemaName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductController) CreateProduct(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	if tenantContext == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant context is missing"})
	}
	return c.SendString("Product created successfully")
}

func (h *ProductController) GetProductByID(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	if tenantContext == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant context is missing"})
	}

	return c.JSON("Product details")
}

func (h *ProductController) UpdateProduct(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	if tenantContext == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant context is missing"})
	}
	return c.SendString("Product updated successfully")
}

func (h *ProductController) DeleteProduct(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	if tenantContext == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant context is missing"})
	}
	return c.SendString("Product deleted successfully")
}
