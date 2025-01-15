package api

import (
	"app/src/api/blog"
	"app/src/api/tag"
	"app/src/api/user"
	"github.com/gofiber/fiber/v3"
)

func Register(s *fiber.App) {
	group := s.Group("/api")
	blog.Register(group)
	tag.Register(group)
	user.Register(group)
}
