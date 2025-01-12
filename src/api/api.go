package api

import (
	"github.com/emrekentli/multitenant-boilerplate/src/api/blog"
	"github.com/emrekentli/multitenant-boilerplate/src/api/tag"
	"github.com/gofiber/fiber/v3"
)

func Register(s *fiber.App) {
	group := s.Group("/api")
	blog.Register(group)
	tag.Register(group)
}
