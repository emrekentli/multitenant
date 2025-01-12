package routes

import (
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/emrekentli/multitenant-boilerplate/src/rest/middlewares"
	"github.com/gofiber/fiber/v2"
	"time"
)

func ApiRoutes(api fiber.Router) {
	api.Get("/ping", Pong)
}
func Pong(c *fiber.Ctx) error {
	tenantContext, ok := c.Locals("tenant").(*middlewares.TenantContext)
	if !ok || tenantContext == nil {
		return c.Status(500).SendString("Tenant context bulunamadı")
	}

	schemaName := tenantContext.SchemaName

	if app.Http.Database.DB == nil {
		return c.Status(500).SendString("Veritabanı bağlantısı bulunamadı")
	}

	query := "SELECT id, name, price, description, created_at FROM " + schemaName + ".products"
	rows, err := app.Http.Database.DB.Query(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString("Veri çekilemedi: " + err.Error())
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var price float64
		var description *string
		var createdAt time.Time

		err = rows.Scan(&id, &name, &price, &description, &createdAt)
		if err != nil {
			return c.Status(500).SendString("Veri okuma hatası: " + err.Error())
		}

		descriptionValue := ""
		if description != nil {
			descriptionValue = *description
		}

		products = append(products, map[string]interface{}{
			"id":          id,
			"name":        name,
			"price":       price,
			"description": descriptionValue,
			"createdAt":   createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	if rows.Err() != nil {
		return c.Status(500).SendString("Sorgu hatası: " + rows.Err().Error())
	}

	return c.JSON(products)
}
