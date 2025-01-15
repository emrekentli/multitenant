package tenantcontext

import (
	"app/src/general/database"
	"context"
	"github.com/gofiber/fiber/v3"
	"log"
)

type TenantContext struct {
	Domain     string
	SchemaName string
}

func RegisterTenantContext(c fiber.Ctx) error {
	tenantDomain := c.Get("X-Tenant-Domain")
	if tenantDomain == "" {
		log.Println("Tenant Domain bulunamadı")
		return fiber.NewError(fiber.StatusBadRequest, "Tenant Domain bulunamadı")
	}

	var schemaName string
	err := database.DB.QueryRow(context.Background(), "SELECT schema_name FROM public.tenants WHERE domain = $1", tenantDomain).Scan(&schemaName)
	if err != nil {
		log.Println("Tenant bulunamadı:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Tenant bulunamadı")
	}

	if schemaName == "" {
		log.Println("Tenant bulunamadı")
		return fiber.NewError(fiber.StatusBadRequest, "Tenant bulunamadı")
	}

	c.Locals("tenant", &TenantContext{
		Domain:     tenantDomain,
		SchemaName: schemaName,
	})

	return c.Next()
}

func GetTenantSchemaName(c fiber.Ctx) string {
	return c.Locals("tenant").(*TenantContext).SchemaName
}
