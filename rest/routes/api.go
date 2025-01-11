package routes

import (
	"database/sql"
	fmt2 "fmt"
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/emrekentli/multitenant-boilerplate/rest/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(api fiber.Router) {
	api.Get("/ping", Pong)
}

func Pong(c *fiber.Ctx) error {
	tenantContext := c.Locals("tenant").(*middlewares.TenantContext)
	tenantID := tenantContext.TenantID

	rows, err := app.Http.Database.DB.Query(fmt2.Sprintf("SELECT * FROM %s.products", tenantID))
	if err != nil {
		return c.Status(500).SendString("Veri çekilemedi")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var products []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var description string
		var created_at string
		var price float64
		err = rows.Scan(&id, &name, &price, &description, &created_at)
		if err != nil {
			return c.Status(500).SendString("Veri okuma hatası")
		}
		products = append(products, map[string]interface{}{
			"id":          id,
			"name":        name,
			"price":       price,
			"description": description,
			"created_at":  created_at,
		})
	}

	if rows.Err() != nil {
		return c.Status(500).SendString("Sorgu hatası")
	}

	// JSON formatında ürün verilerini döndür
	return c.JSON(products)
}
