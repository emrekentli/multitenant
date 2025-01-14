package middlewares

import (
	"app/app/middlewares/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"os"
)

func RegisterBuiltInMiddlewares(s *fiber.App) {
	s.Use(etag.New())
	s.Use(helmet.New())
	s.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	if os.Getenv("ENABLE_LOGGER") == "true" {
		s.Use(logger.New())
	}
	s.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
	jwt.RegisterJwtMiddleware(s)
	s.Get("/*", static.New("./public"))
}
