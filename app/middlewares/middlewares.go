package middlewares

import (
	"app/app/middlewares/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func RegisterBuiltInMiddlewares(s *fiber.App) {
	s.Use(etag.New())

	s.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
	s.Use(compress.New())
	jwt.RegisterJwtMiddleware(s)
	s.Get("/*", static.New("./public"))
}
