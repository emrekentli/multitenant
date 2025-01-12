package middlewares

import (
	"context"
	"github.com/emrekentli/multitenant-boilerplate/config/database"
	"github.com/gofiber/fiber/v3"
	"log"
)

type TenantContext struct {
	Domain     string
	SchemaName string
}

func RegisterTenantContext(app *fiber.App) {
	app.Use(
		func(c fiber.Ctx) error {
			tenantDomain := c.Get("X-Tenant-Domain")
			if tenantDomain == "" {
				log.Println("Tenant Domain bulunamadı")
				return c.Status(400).SendString("Tenant Domain bulunamadı")
			}
			var schemaName string
			err := database.DB.QueryRow(context.Background(), "SELECT schema_name FROM public.tenants WHERE domain = $1", tenantDomain).Scan(&schemaName)
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
		})
}
