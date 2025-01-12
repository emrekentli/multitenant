package middlewares

import (
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/gofiber/fiber/v2"
	"log"
)

type TenantContext struct {
	Domain     string
	SchemaName string
}

func SetTenantContext(c *fiber.Ctx) error {
	tenantDomain := c.Get("X-Tenant-Domain")
	if tenantDomain == "" {
		log.Println("Tenant Domain bulunamadı")
		return c.Status(400).SendString("Tenant Domain bulunamadı")
	}

	var schemaName string
	err := app.Http.Database.DB.QueryRow(c.Context(), "SELECT schema_name FROM public.tenants WHERE domain = $1", tenantDomain).Scan(&schemaName)
	if err != nil {
		log.Println("Tenant bulunamadı:", err)
		return c.Status(400).SendString("Tenant bulunamadı")
	}

	if schemaName == "" {
		log.Println("Tenant bulunamadı")
		return c.Status(400).SendString("Tenant bulunamadı")
	}

	c.Locals("tenant", &TenantContext{
		Domain:     tenantDomain,
		SchemaName: schemaName,
	})

	return c.Next()
}
