package jwt

import (
	"app/config"
	"app/src/general/util/rest"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strings"
	"time"
)

var bypassJwtRoutes = []string{
	"POST /api/user/login",
	"GET /api/blog",
	"GET /public/*",
	"GET /healthcheck",
}

func pathMatches(pattern, path string) bool {
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(path, prefix)
	}
	return pattern == path
}

func filterJwtMiddleware(c fiber.Ctx) bool {
	for _, entry := range bypassJwtRoutes {
		parts := strings.SplitN(entry, " ", 2)
		if len(parts) != 2 {
			continue
		}
		method, pattern := parts[0], parts[1]
		if c.Method() == method && pathMatches(pattern, c.Path()) {
			return true
		}
	}
	return false
}

func RegisterJwtMiddleware(app *fiber.App) {
	app.Use(func(c fiber.Ctx) error {
		if filterJwtMiddleware(c) {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			authHeader = c.Cookies("Authorization")
		}
		if strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}
		authHeader = strings.TrimSpace(authHeader)

		if authHeader == "" {
			return onUnAuthorized(c)
		}

		token, err := VerifyToken(authHeader)
		if err != nil || !token.Valid {
			log.Printf("JWT parse error: %v", err)
			return onUnAuthorized(c)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return onUnAuthorized(c)
		}

		// Bilgileri sonraki handler'lara ilet
		c.Locals("ID", claims["ID"])
		c.Locals("permissions", claims["permissions"])
		return c.Next()
	})
}

func onUnAuthorized(c fiber.Ctx) error {
	return rest.ErrorRes(c, rest.Unauthorized, rest.ErrorCode[rest.Unauthorized])
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	// Exp kontrolü (çoğu zaman gerekmez ama güvenlik için ekli)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Token expired")
			}
		}
	}
	return token, nil
}

func CreateJwt(id int64, permissions []*string) (string, error) {
	day := 24 * time.Hour
	claims := jwt.MapClaims{
		"ID":          id,
		"exp":         time.Now().Add(day).Unix(),
		"permissions": permissions,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(config.JwtSecretByte)
	if err != nil {
		return "", err
	}
	return t, nil
}

func RequirePermission(permission string) fiber.Handler {
	fmt.Println("RequirePermission middleware registered for permission:", permission)
	return func(c fiber.Ctx) error {
		perms, ok := c.Locals("permissions").([]interface{})
		if !ok {
			return rest.ErrorRes(c, rest.Unauthorized, "Yetersiz yetki")
		}
		for _, p := range perms {
			if ps, ok := p.(string); ok && ps == permission {
				return c.Next()
			}
		}
		return rest.ErrorRes(c, rest.Unauthorized, "Yetersiz yetki")
	}
}
