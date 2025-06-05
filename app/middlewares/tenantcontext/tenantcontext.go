package tenantcontext

import (
	"app/src/general/cache"
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

	schemaName, exists := cache.GetSchemaByDomain(tenantDomain)
	if !exists {
		log.Printf("Tenant Domain '%s' için şema bulunamadı", tenantDomain)
		return fiber.NewError(fiber.StatusNotFound, "Tenant Domain için şema bulunamadı")
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
