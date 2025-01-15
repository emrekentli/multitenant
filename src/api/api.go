package api

import (
	"app/app/middlewares/tenantcontext"
	"app/src/api/blog"
	"app/src/api/role"
	"app/src/api/tag"
	"app/src/api/user"
	"github.com/gofiber/fiber/v3"
)

func Register(s *fiber.App) {
	group := s.Group("/api")
	group.Use(tenantcontext.RegisterTenantContext)
	// Auth
	role.Register(group)
	user.Register(group)

	// Other
	blog.Register(group)
	tag.Register(group)
}
