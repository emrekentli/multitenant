package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type TenantContext struct {
	TenantID string
}

func SetTenantContext(c *fiber.Ctx) error {
	// get from header
	tenantID := c.Get("X-Tenant-ID")

	if tenantID == "" {
		// Hata durumunda ya da varsayılan bir tenant kimliği kullanılabilir.
		log.Println("Tenant ID bulunamadı")
		tenantID = "public"
	}

	c.Locals("tenant", &TenantContext{TenantID: tenantID})

	return c.Next()
}
