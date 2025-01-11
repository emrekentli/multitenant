package routes

import (
	"database/sql"
	"fmt"
	"github.com/emrekentli/multitenant-boilerplate/app"
	"github.com/emrekentli/multitenant-boilerplate/rest/middlewares"
	"github.com/gofiber/fiber/v2"
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

	query := fmt.Sprintf("SELECT * FROM %s.products", schemaName)
	rows, err := app.Http.Database.DB.Query(query)
	if err != nil {
		return c.Status(500).SendString("Veri çekilemedi: " + err.Error())
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
		var createdAt string
		var price float64
		err = rows.Scan(&id, &name, &price, &description, &createdAt)
		if err != nil {
			return c.Status(500).SendString("Veri okuma hatası: " + err.Error())
		}
		products = append(products, map[string]interface{}{
			"id":          id,
			"name":        name,
			"price":       price,
			"description": description,
			"createdAt":   createdAt,
		})
	}

	if rows.Err() != nil {
		return c.Status(500).SendString("Sorgu hatası: " + rows.Err().Error())
	}

	return c.JSON(products)
}
